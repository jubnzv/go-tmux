// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

type Window struct {
	Name           string
	Id             int
	SessionId      int
	SessionName    string
	StartDirectory string // Path to window working directory
	Panes          []Pane // List of panes used in initial window configuration
}

// Return list with all panes for this window.
func (w *Window) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-s", "-t", w.Name})
}

// Add pane to window configuration. This will change only in-library
// window representation. Used for initial configuration before creating new
// window.
func (w *Window) AddPane(pane Pane) {
	w.Panes = append(w.Panes, pane)
}
