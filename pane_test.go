// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

import (
	"os"
	"testing"
)

func TestPaneGetCurrentPath(t *testing.T) {
	if InTravis() {
		t.Skip("Skipping this test in travis.")
	}

	err := os.Chdir("/tmp")
	if err != nil {
		t.Errorf("There are some problems with /tmp directory: %s", err)
	}

	s := createSession()
	restoreSession()
	s.AttachSession()
	defer restoreSession()
	defer sessionsReaper(s.Name)
	window, _ := s.NewWindow("test-window")
	panes, _ := window.ListPanes()
	pane := panes[0]

	path, err := pane.GetCurrentPath()
	if path != "/tmp" {
		t.Errorf("Incorrect path (expected %s got %s)", "/tmp", path)
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}
