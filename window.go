// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

package tmux

const (
	LayoutEvenHorizontal = "even-horizontal"
	LayoutEvenVertical   = "even-vertical"
	LayoutMainHorizontal = "main-horizontal"
	LayoutMainVertical   = "main-vertical"
	LayoutTiled          = "tiled"
)

// Represents a tmux window:
// https://github.com/tmux/tmux/wiki/Getting-Started#sessions-windows-and-panes
type Window struct {
	Name           string
	Id             int
	SessionId      int
	SessionName    string
	StartDirectory string // Path to window working directory
	Layout         string // Preset arrangements of panes
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
		Layout:         "",
	}
}

// Returns a list with all panes for this window.
func (w *Window) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-s", "-t", w.Name})
}

// Adds the pane to the window configuration. This will change only in-library
// window representation. Used for initial configuration before creating a new
// window.
func (w *Window) AddPane(pane Pane) {
	w.Panes = append(w.Panes, pane)
}

// Sets the window layout. The possible value can be one of the contants:
// LayoutEvenVertical, LayoutEvenHorizontal, LayoutMainVertical,  LayoutMainHorizontal, LayoutTiled
// See: https://www.man7.org/linux/man-pages/man1/tmux.1.html#WINDOWS_AND_PANES
func (w *Window) SetLayout(layout string) {
	w.Layout = layout
}
