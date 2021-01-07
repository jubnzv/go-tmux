// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import (
	"testing"
)

func TestWindowListPanes(t *testing.T) {
	if InTravis() {
		t.Skip("Skipping this test in travis.")
	}

	s := createSession()
	restoreSession()
	s.AttachSession()
	defer restoreSession()
	defer sessionsReaper(s.Name)
	w, _ := s.NewWindow("test-window")
	panes, _ := w.ListPanes()

	for _, p := range panes {
		if p.SessionId != s.Id {
			t.Fatalf("Incorrect session id (expected %d got %d)", s.Id, p.SessionId)
		}
		if p.SessionName != s.Name {
			t.Fatalf("Incorrect session name (expected %s got %s)", s.Name, p.SessionName)
		}
		if p.WindowId != w.Id {
			t.Fatalf("Incorrect window id (expected %d got %d)", w.Id, p.WindowId)
		}
		if p.WindowName != w.Name {
			t.Fatalf("Incorrect window name (expected %s got %s)", w.Name, p.WindowName)
		}
	}
}

func TestWindowHaveSinglePaneAfterInit(t *testing.T) {
	if InTravis() {
		t.Skip("Skipping this test in travis.")
	}

	s := createSession()
	restoreSession()
	s.AttachSession()
	defer restoreSession()
	defer sessionsReaper(s.Name)
	w, _ := s.NewWindow("test-window")
	panes, err := w.ListPanes()
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(panes) != 1 {
		t.Fatalf("Window must have single pane after init (got %d)", len(panes))
	}
}
