// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Represent a tmux pane:
// https://github.com/tmux/tmux/wiki/Getting-Started#sessions-windows-and-panes
type Pane struct {
	ID          int
	SessionId   int
	SessionName string
	WindowId    int
	WindowName  string
	WindowIndex int
	Active      bool
	Index       int
}

// Creates a new pane object.
func NewPane(id int, sessionId int, sessionName string, windowId int,
	windowName string, windowIndex int, active bool,
) *Pane {
	return &Pane{
		ID:          id,
		SessionId:   sessionId,
		SessionName: sessionName,
		WindowId:    windowId,
		WindowName:  windowName,
		WindowIndex: windowIndex,
		Active:      active,
	}
}

// Return a list of panes. Optional arguments are define the search scope with
// tmux command keys (see tmux(1) manpage):
//
// list-panes [-as] [-F format] [-t target]
//   - `-a`: target is ignored and all panes on the server are listed
//   - `-s`: target is a session. If neither is given, target is a window (or
//     the current window).
func ListPanes(args []string) ([]Pane, error) {
	format := strings.Join([]string{
		"#{session_id}",
		"#{session_name}",
		"#{window_id}",
		"#{window_name}",
		"#{window_index}",
		"#{pane_id}",
		"#{pane_active}",
		"#{pane_index}",
	}, ":")

	args = append([]string{"list-panes", "-F", format}, args...)

	out, _, err := RunCmd(args)
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(out, "\n")
	panes := []Pane{}
	re := regexp.MustCompile(`\$([0-9]+):(.+):@([0-9]+):(.+):([0-9]+):%([0-9]+):([01]):([0-9]+)`)
	const paneParts = 8

	for _, line := range outLines {
		result := re.FindStringSubmatch(line)
		if len(result) <= paneParts {
			continue
		}

		sessionID, errAtoi := strconv.Atoi(result[1])
		if errAtoi != nil {
			return nil, errAtoi
		}

		windowID, errAtoi := strconv.Atoi(result[3])
		if errAtoi != nil {
			return nil, errAtoi
		}

		windowIndex, errAtoi := strconv.Atoi(result[5])
		if errAtoi != nil {
			return nil, errAtoi
		}

		paneId, errAtoi := strconv.Atoi(result[6])
		if errAtoi != nil {
			return nil, errAtoi
		}

		paneIndex, errAtoi := strconv.Atoi(result[8])
		if errAtoi != nil {
			return nil, errAtoi
		}

		panes = append(panes, Pane{
			ID:          paneId,
			SessionId:   sessionID,
			SessionName: result[2],
			WindowId:    windowID,
			WindowName:  result[4],
			WindowIndex: windowIndex,
			Active:      result[7] == "1",
			Index:       paneIndex,
		})
	}

	return panes, nil
}

// Returns current path for this pane.
func (p *Pane) GetCurrentPath() (string, error) {
	args := []string{
		"display-message",
		"-P", "-F", "#{pane_current_path}",
	}
	out, _, err := RunCmd(args)
	if err != nil {
		return "", err
	}

	// Remove trailing CR
	out = out[:len(out)-1]

	return out, nil
}

// Log the pane buffer to a file
func (p *Pane) Pipe(path string) (string, error) {
	args := []string{
		"pipe-pane",
		"-t", fmt.Sprintf("%s:%s.%d", p.SessionName, p.WindowName, p.Index),
		"-o", fmt.Sprintf("cat >>%s", path),
	}

	out, stdErr, err := RunCmd(args)
	if err != nil {
		return stdErr, err
	}

	return out, nil
}

// Capture the pane buffer
func (p *Pane) Capture() (string, error) {
	args := []string{
		"capture-pane",
		"-t", fmt.Sprintf("%s:%s.%d", p.SessionName, p.WindowName, p.Index),
		"-p",
	}

	out, stdErr, err := RunCmd(args)
	if err != nil {
		return stdErr, err
	}

	// Do not remove the tailing CR,
	// maybe it's important for the caller
	// for capture-pane.
	return out, nil
}

// RunCommand runs a command in the pane.
func (p *Pane) RunCommand(command string) error {
	args := []string{
		"send-keys",
		"-t", fmt.Sprintf("%s:%s.%d", p.SessionName, p.WindowName, p.Index),
		command,
		"C-m",
	}
	_, stdErr, err := RunCmd(args)
	if err != nil {
		return fmt.Errorf("%v: %s", err, stdErr)
	}
	return nil
}

// Selects the pane.
func (p *Pane) Select() error {
	args := []string{
		"select-pane",
		"-t", fmt.Sprintf("%s:%s.%d", p.SessionName, p.WindowName, p.Index),
	}
	_, stdErr, err := RunCmd(args)
	if err != nil {
		return fmt.Errorf("%v: %s", err, stdErr)
	}
	return nil
}
