// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

import (
	"testing"
)

func createSession() Session {
	server := new(Server)
	session, _ := server.NewSession("test-session")
	return session
}

func TestListWindows(t *testing.T) {
	s := createSession()
	defer sessionsReaper(s.Name)
	if _, err := s.ListWindows(); err != nil {
		t.Fatalf("ListWindows: %s", err)
	}
}

func TestNewWindow(t *testing.T) {
	s := createSession()
	defer sessionsReaper(s.Name)

	window, err := s.NewWindow("test-new-window")
	if err != nil {
		t.Fatalf("NewWindow: %s", err)
	}

	ws, _ := s.ListWindows()

	// Check created window name
	found := false
	for _, w := range ws {
		if w.Name == window.Name {
			found = true
			break
		}
	}
	if found == false {
		t.Fatalf("Can't find created window 'test-new-window'")
	}

	// Check created window id
	found = false
	for _, w := range ws {
		if w.Id == window.Id {
			found = true
			break
		}
	}
	if found == false {
		t.Fatalf("Can't find created window by id: %d", window.Id)
	}

	if len(window.SessionName) == 0 {
		t.Fatalf("New window created in inappropriate session (expected %s got %s)", s.Name, window.SessionName)
	}
	if window.SessionName != s.Name {
		t.Fatalf("New window created in inappropriate session (expected %s got %s)", s.Name, window.SessionName)
	}
	if window.SessionId != s.Id {
		t.Fatalf("New window: incorrect session id (expected %d, got %d)", s.Id, window.SessionId)
	}
}

func TestSessionListPanes(t *testing.T) {
	s := createSession()
	defer sessionsReaper(s.Name)
	panes, _ := s.ListPanes()

	for _, p := range panes {
		if p.SessionId != s.Id {
			t.Fatalf("Incorrect session id (expected %d got %d)", s.Id, p.SessionId)
		}
		if p.SessionName != s.Name {
			t.Fatalf("Incorrect session name (expected %s got %s)", s.Name, p.SessionName)
		}
	}
}

func TestGetAttachedSessionName(t *testing.T) {
	s := createSession()
	defer sessionsReaper(s.Name)

	name, err := GetAttachedSessionName()
	if err != nil {
		t.Fatalf("GetAttachedSessionName: %s", err)
	}

	// Skipped in local testing because I use it inside tmux
	if InTravis() {
		if name != s.Name {
			t.Fatalf("Incorrect session name (expected %s got %s)", s.Name, name)
		}
	}
}
