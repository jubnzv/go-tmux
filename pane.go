// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	paneParts = 7
)

type Pane struct {
	ID          int
	SessionId   int
	SessionName string
	WindowId    int
	WindowName  string
	WindowIndex int
	Active      bool
}

// Return list of panes. Optional arguments are define the search scope with
// tmux command keys (see tmux(1) manpage):
// list-panes [-as] [-F format] [-t target]
//
// * `-a`: target is ignored and all panes on the server are listed
// * `-s`: target is a session. If neither is given, target is a window (or
//   the current window).
func ListPanes(args []string) ([]Pane, error) {
	format := strings.Join([]string{
		"#{session_id}",
		"#{session_name}",
		"#{window_id}",
		"#{window_name}",
		"#{window_index}",
		"#{pane_id}",
		"#{pane_active}",
	}, ":")

	args = append([]string{"list-panes", "-F", format}, args...)

	out, _, err := RunCmd(args)
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(out, "\n")
	panes := []Pane{}
	re := regexp.MustCompile(`\$([0-9]+):(.+):@([0-9]+):(.+):([0-9]+):%([0-9]+):([01])`)

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

		paneIndex, errAtoi := strconv.Atoi(result[6])
		if errAtoi != nil {
			return nil, errAtoi
		}

		panes = append(panes, Pane{
			SessionId:   sessionID,
			SessionName: result[2],
			WindowId:    windowID,
			WindowName:  result[4],
			WindowIndex: windowIndex,
			ID:          paneIndex,
			Active:      result[7] == "1",
		})
	}

	return panes, nil
}

// Returns current path for this pane.
func (p *Pane) GetCurrentPath() (string, error) {
	args := []string{
		"display-message",
		"-P", "-F", "#{pane_current_path}"}
	out, _, err := RunCmd(args)
	if err != nil {
		return "", err
	}

	// Remove trailing CR
	out = out[:len(out)-1]

	return out, nil
}

func (p *Pane) Capture() (string, error) {
	args := []string{
		"capture-pane",
		"-t",
		fmt.Sprintf("%%%d", p.ID),
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
