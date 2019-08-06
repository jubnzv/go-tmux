// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

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

func InTravis() bool {
	if os.Getenv("IN_TRAVIS") == "1" {
		return true
	} else {
		return false
	}
}
