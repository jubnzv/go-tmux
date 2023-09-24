// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

// tmux configuration used to setup workflow with user-defined sessions,
// windows and panes.

package tmux

import (
	"errors"
	"fmt"
	"strings"
)

type Configuration struct {
	Server        *Server    // Pointer to used tmux server
	Sessions      []*Session // List of sessions to be initialized
	ActiveSession *Session   // Session to be attached after initialization.
	// If nil, leave sessions detached.
}

// Checks that given configuration is correct
func (c *Configuration) checkInput() error {
	for _, s := range c.Sessions {
		// xxx: What is upper limit for tmux windows number?..
		if len(s.Windows) == 0 {
			msg := fmt.Sprintf("Session %s doesn't contain any windows!", s.Name)
			return errors.New(msg)
		}
	}

	return nil
}

// Apply given configuration to setup a user-defined workflow
// Before running this method, user must make sure that there is no windows and
// session with same names exists. Otherwise existing sessions/windows will be
// replaced with the new ones.
func (c *Configuration) Apply() error {
	if c.Server == nil {
		return errors.New("Server was not initialized")
	}
	if len(c.Sessions) == 0 {
		return errors.New("Requiered at least single tmux session to apply configuration")
	}

	// Check for requested configuration correctness
	if err := c.checkInput(); err != nil {
		return err
	}

	// Initialize sessions
	for _, s := range c.Sessions {
		// Set initial window for a new session
		initial_window := s.Windows[0]

		// Select start directory for a session
		args_start_dir := []string{}
		if len(s.StartDirectory) != 0 {
			args_start_dir = []string{"-c", s.StartDirectory}
		}

		// Should this session be attached after init?
		attached_key := "-d"
		if s == c.ActiveSession {
			attached_key = ""
		}

		// Start a new session
		args := []string{
			"new-session",
			attached_key,
			"-n", initial_window.Name,
			"-D", // If session with same name exists, attach to it
			"-s", s.Name,
		}
		args = append(args, args_start_dir...)
		_, err_out, err_exec := RunCmd(args)
		if err_exec != nil {
			// It's okay, if session already exists.
			if !strings.Contains(err_out, "exit status 1") {
				return err_exec
			}
		}

		// Add windows for created session
		for _, w := range s.Windows {
			// Select start directory for this window
			// If empty, use StartDirectory from session
			var windowStartDirectory string
			if len(w.StartDirectory) != 0 {
				windowStartDirectory = w.StartDirectory
			} else if len(s.StartDirectory) != 0 {
				windowStartDirectory = s.StartDirectory
			}
			winId := fmt.Sprintf("%s:%d", s.Name, w.Id)
			args_start_dir := []string{}

			if len(windowStartDirectory) != 0 {
				args_start_dir = []string{"-c", windowStartDirectory}
			}

			// Create a new window
			args = []string{
				"new-window",
				"-k", // Destroy windows if already exists
				"-n", w.Name,
				"-t", winId,
			}
			args = append(args, args_start_dir...)
			_, _, err_exec := RunCmd(args)
			if err_exec != nil {
				return err_exec
			}

			// Setup panes for created window
			for idx := range w.Panes {
				// First pane is created automatically, so split existing window
				if idx != 0 {
					args = []string{
						"split-window",
						"-t", winId,
						"-c", windowStartDirectory,
					}
					_, _, err_exec := RunCmd(args)
					if err_exec != nil {
						return err_exec
					}
				}
			}

			// Select layout if defined
			if len(w.Layout) != 0 {
				args = []string{"select-layout", "-t", winId, w.Layout}
				_, _, err_exec := RunCmd(args)
				if err_exec != nil {
					return err_exec
				}
			}
		}
	}

	return nil
}
