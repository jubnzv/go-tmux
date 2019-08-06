// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import (
	"testing"
    "os"
)

func TestPaneGetCurrentPath(t *testing.T) {
    err := os.Chdir("/tmp")
    if err != nil {
        t.Errorf("There are some problems with /tmp directory: %s", err)
    }

	s := createSession()
	s.AttachSession()
	defer s.DettachSession()
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
