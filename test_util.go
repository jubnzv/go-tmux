// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

import (
	"os"
	"strings"
)

// Kills sessions that contains namePattern substring in the name.
func sessionsReaper(namePattern string) {
	s := new(Server)
	// Suppose that ListSession works.
	sessions, _ := s.ListSessions()
	for _, ss := range sessions {
		if strings.Contains(ss.Name, namePattern) {
			s.KillSession(ss.Name)
		}
	}
}

// Restores the session that was active before test.
func restoreSession() {
	if !IsInsideTmux() {
		return
	}

	// Session name when test is running
	session_name := ""

	restore_fun := func() {
		if session_name == "" {
			session_name, _ = GetAttachedSessionName()
		} else {
			server := new(Server)
			sessions, _ := server.ListSessions()
			for _, s := range sessions {
				if s.Name == session_name {
					s.AttachSession()
				}
			}
		}
		return
	}

	restore_fun()
}

func InTravis() bool {
	if os.Getenv("IN_TRAVIS") == "1" {
		return true
	} else {
		return false
	}
}
