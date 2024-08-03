// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Represents a tmux session:
// https://github.com/tmux/tmux/wiki/Getting-Started#sessions-windows-and-panes
type Session struct {
	Id             int      // Session id
	Name           string   // Session name
	StartDirectory string   // Path to window start directory
	Windows        []Window // List of windows used on session initialization
}

// Creates a new session object.
func NewSession(id int, name, startDirectory string, windows []Window) *Session {
	return &Session{
		Id:             id,
		Name:           name,
		StartDirectory: startDirectory,
		Windows:        windows,
	}
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

// Adds the window to the session configuration. This will change only
// in-library session representation. Used for initial configuration before
// creating a new session.
func (s *Session) AddWindow(window Window) {
	s.Windows = append(s.Windows, window)
}

// Lists all windows in the current session.
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
			Id:             id,
			StartDirectory: result[3],
			SessionName:    s.Name,
			SessionId:      s.Id})
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

// Detaches from the current session.
// Detaching from the tmux session means that the client exits and detaches
// from the outside terminal.
// See: https://github.com/tmux/tmux/wiki/Getting-Started#attaching-and-detaching
func (s *Session) DettachSession() error {
	args := []string{
		"detach-client",
		"-s", s.Name}
	if err := ExecCmd(args); err != nil {
		return err
	}
	return nil
}

// Creates a new window inside this session.
func (s *Session) NewWindow(name string, args ...string) (window Window, err error) {
	args = append([]string{
		"new-window",
		"-d",
		"-t", fmt.Sprintf("%s:", s.Name),
		"-n", name,
		"-F", "#{window_id}:#{window_name}", "-P"},
		args...)

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
		SessionId:   s.Id,
		SessionName: s.Name,
		WindowId:    id,
		WindowName:  result[2],
		WindowIndex: 0}
	new_window := Window{
		Name:        result[2],
		Id:          id,
		SessionName: s.Name,
		SessionId:   s.Id,
		Panes:       []Pane{pane}}

	// Append the Window struct to the Windows slice in the Session struct.
	s.Windows = append(s.Windows, new_window)

	return new_window, nil
}

// Returns list with all panes for this session.
func (s *Session) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-s", "-t", s.Name})
}

// Returns a name of the attached tmux session.
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
