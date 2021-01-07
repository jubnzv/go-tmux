// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Session struct {
	Name           string   // Session name
	Id             int      // Session id
	StartDirectory string   // Path to window start directory
	Windows        []Window // List of windows used on session initialization
}

// Checks tmux rules for sessions naming. Reference:
// https://github.com/tmux/tmux/blob/5489796737108cb9bba01f831421e531a50b946b/session.c#L238
func checkSessionName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if strings.Contains(name, ":") {
		return false
	}
	if strings.Contains(name, ".") {
		return false
	}
	return true
}

// Add window to session configuration. This will change only in-library
// session representation. Used for initial configuration before creating new
// session.
func (s *Session) AddWindow(window Window) {
	s.Windows = append(s.Windows, window)
}

func (s *Session) InitSession() error {
	if len(s.Windows) == 0 {
		return errors.New("Nothing to do.")
	}

	return nil
}

// List all windows in the current session.
func (s *Session) ListWindows() ([]Window, error) {
	args := []string{
		"list-windows",
		"-t", s.Name,
		"-F", "#{window_id}:#{window_name}:#{pane_current_path}"}

	out, _, err := RunCmd(args)
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(out, "\n")
	windows := []Window{}
	re := regexp.MustCompile(`@([0-9]+):(.+):(.+)`)
	for _, line := range outLines {
		result := re.FindStringSubmatch(line)
		if len(result) < 4 {
			continue
		}
		id, err_atoi := strconv.Atoi(result[1])
		if err_atoi != nil {
			return nil, err_atoi
		}

		windows = append(windows, Window{
			Name:           result[2],
			ID:             id,
			StartDirectory: result[3],
			SessionName:    s.Name,
			SessionID:      s.Id})
	}

	return windows, nil
}

// Attach to existing tmux session.
func (s *Session) AttachSession() error {
	args := []string{}
	// If run inside tmux, switch the current session to the new one.
	if !IsInsideTmux() {
		args = append(args, "attach-session", "-t", s.Name)
	} else {
		args = append(args, "switch-client", "-t", s.Name)
	}

	if err := ExecCmd(args); err != nil {
		return err
	}

	return nil
}

// Detach from current session.
func (s *Session) DettachSession() error {
	args := []string{
		"detach-client",
		"-s", s.Name}
	if err := ExecCmd(args); err != nil {
		return err
	}
	return nil
}

// Create a new window.
func (s *Session) NewWindow(name string) (window Window, err error) {
	args := []string{
		"new-window",
		"-d",
		"-t", fmt.Sprintf("%s:", s.Name),
		"-n", name,
		"-F", "#{window_id}:#{window_name}", "-P"}
	out, _, err_exec := RunCmd(args)
	if err_exec != nil {
		return window, err_exec
	}

	re := regexp.MustCompile(`@([0-9]+):(.+)`)
	result := re.FindStringSubmatch(out)
	if len(result) < 3 {
		return window, errors.New("Error creating new window")
	}
	id, err_atoi := strconv.Atoi(result[1])
	if err_atoi != nil {
		return window, err_atoi
	}

	pane := Pane{
		SessionID:   s.Id,
		SessionName: s.Name,
		WindowID:    id,
		WindowName:  result[2],
		WindowIndex: 0}
	new_window := Window{
		Name:        result[2],
		ID:          id,
		SessionName: s.Name,
		SessionID:   s.Id,
		Panes:       []Pane{pane}}
	return new_window, nil
}

// Return list with all panes for this session.
func (s *Session) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-s", "-t", s.Name})
}

// Return name of attached tmux session.
func GetAttachedSessionName() (string, error) {
	args := []string{
		"display-message",
		"-p", "#S"}
	out, _, err := RunCmd(args)
	if err != nil {
		return "", err
	}

	// Remove trailing CR
	out = out[:len(out)-1]

	return out, nil
}
