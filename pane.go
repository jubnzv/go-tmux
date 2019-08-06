// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import (
	"regexp"
    "fmt"
	"strconv"
	"strings"
)

type Pane struct {
	SessionId   int
	SessionName string
	WindowId    int
	WindowName  string
	WindowIndex int
}

// Return list of panes. Optional arguments are define the search scope with
// tmux command keys (see tmux(1) manpage):
// list-panes [-as] [-F format] [-t target]
//
// * `-a`: target is ignored and all panes on the server are listed
// * `-s`: target is a session. If neither is given, target is a window (or
//   the current window).
func ListPanes(args []string) ([]Pane, error) {
	args = append([]string{"list-panes", "-F", "#{session_id}:#{session_name}:#{window_id}:#{window_name}:#{window_index}"}, args...)

	out, _, err := RunCmd(args)
	if err != nil {
        fmt.Printf("RUN ERR: %s", err)
		return nil, err
	}

	outLines := strings.Split(out, "\n")
	panes := []Pane{}
	re := regexp.MustCompile(`\$([0-9]+):(.+):@([0-9]+):(.+):([0-9]+)`)
	for _, line := range outLines {
		result := re.FindStringSubmatch(line)
		if len(result) < 6 {
			continue
		}
		session_id, err_atoi := strconv.Atoi(result[1])
		if err_atoi != nil {
			return nil, err_atoi
		}
		window_id, err_atoi := strconv.Atoi(result[3])
		if err_atoi != nil {
			return nil, err_atoi
		}
		window_index, err_atoi := strconv.Atoi(result[5])
		if err_atoi != nil {
			return nil, err_atoi
		}

		panes = append(panes, Pane{
			SessionId:   session_id,
			SessionName: result[2],
			WindowId:    window_id,
			WindowName:  result[4],
			WindowIndex: window_index})
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
