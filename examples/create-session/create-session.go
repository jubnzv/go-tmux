// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

// Simple example that shows how to create a tmux session with user-defined
// configuration.

package main

import (
	"fmt"

	gotmux "github.com/jubnzv/go-tmux"
)

func main() {
	// Create instance of running tmux server.
	server := new(gotmux.Server)

	// Check that "example" session already exists.
	exists, err := server.HasSession("example")
	if err != nil {
		msg := fmt.Errorf("Can't check 'example' session: %s", err)
		fmt.Println(msg)
		return
	}
	if exists {
		// Sure, you can use KillSession here.
		fmt.Println("Session 'example' already exists!")
		fmt.Println("Please stop it before running this demo.")
		return
	}

	// Prepare configuration for a new session with some windows.
	session := gotmux.Session{Name: "example-session"}
	w1 := gotmux.Window{Name: "first", ID: 0}
	w2 := gotmux.Window{Name: "second", ID: 1}
	session.AddWindow(w1)
	session.AddWindow(w2)
	server.AddSession(session)
	sessions := []*gotmux.Session{}
	sessions = append(sessions, &session)
	conf := gotmux.Configuration{
		Server:        server,
		Sessions:      sessions,
		ActiveSession: nil}

	// Setup this configuration.
	err = conf.Apply()
	if err != nil {
		msg := fmt.Errorf("Can't apply prepared configuration: %s", err)
		fmt.Println(msg)
		return
	}

	// Attach to created session
	err = session.AttachSession()
	if err != nil {
		msg := fmt.Errorf("Can't attached to created session: %s", err)
		fmt.Println(msg)
		return
	}
}
