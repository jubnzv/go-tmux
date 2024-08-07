// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

// tmux configuration used to setup workflow with user-defined sessions,
// windows and panes.

package tmux

import (
	"errors"
	"fmt"
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
	for si, s := range c.Sessions {
		// Set initial window for a new session
		initial_window := s.Windows[0]

		extra_args := []string{"-n", initial_window.Name}

		// Select start directory for a session
		if len(s.StartDirectory) != 0 {
			extra_args = append(extra_args, "-c", s.StartDirectory)
		}

		session, err := c.Server.NewSession(s.Name, extra_args...)

		if err != nil {
			return err
		}

		s.Name = session.Name
		s.Id = session.Id

		if s == c.ActiveSession {
			s.AttachSession()
		}

		// Add windows for created session
		var win []Window
		for _, w := range s.Windows {
			// Select start directory for this window
			// If empty, use StartDirectory from session
			if len(w.StartDirectory) == 0 && len(s.StartDirectory) != 0 {
				w.StartDirectory = s.StartDirectory
			}

			if w.Name == initial_window.Name {
				w.Name = initial_window.Name
				w.Id = 0
			} else {
				extra_args = []string{}

				if len(w.StartDirectory) != 0 {
					extra_args = append(extra_args, "-c", w.StartDirectory)
				}

				// Create a new window
				window, err := s.NewWindow(w.Name, extra_args...)
				if err != nil {
					return err
				}

				w.Name = window.Name
				w.Id = window.Id
			}

			w.SessionName = s.Name
			w.SessionId = s.Id

			// Setup panes for created window
			orig_panes := w.Panes
			w.Panes, _ = w.ListPanes()
			for idx := range orig_panes {
				// First pane is created automatically, so split existing window
				if idx > 0 {
					// Create a new pane
					pane, err := w.SplitPane()
					if err != nil {
						return err
					}
					w.Panes[idx] = pane
				}
			}

			// Select layout if defined
			if len(w.Layout) != 0 {
				args := []string{"select-layout", "-t", fmt.Sprintf("%v", w.Id), w.Layout}
				_, _, err_exec := RunCmd(args)
				if err_exec != nil {
					return err_exec
				}
			}

			win = append(win, w)
		}

		c.Sessions[si].Windows = win
		c.Sessions[si] = s
	}

	return nil
}
