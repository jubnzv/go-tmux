// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

// Represents a tmux window:
// https://github.com/tmux/tmux/wiki/Getting-Started#sessions-windows-and-panes
type Window struct {
	Name           string
	Id             int
	SessionId      int
	SessionName    string
	StartDirectory string // Path to window working directory
	Panes          []Pane // List of panes used in initial window configuration
}

// Creates a new window object.
func NewWindow(id int, name string, sessionId int, sessionName string, startDirectory string, panes []Pane) *Window {
	return &Window{
		Name:           name,
		Id:             id,
		SessionId:      sessionId,
		SessionName:    sessionName,
		StartDirectory: startDirectory,
		Panes:          panes,
	}
}

// Returns a list with all panes for this window.
func (w *Window) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-t", w.Name})
}

// Adds the pane to the window configuration. This will change only in-library
// window representation. Used for initial configuration before creating a new
// window.
func (w *Window) AddPane(pane Pane) {
	w.Panes = append(w.Panes, pane)
}
