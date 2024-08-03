// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>
//
// Represents a tmux server - an object that holds sessions with their windows
// and panes.

package tmux

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Represents a tmux server:
// https://github.com/tmux/tmux/wiki/Getting-Started#the-tmux-server-and-clients
type Server struct {
	SocketPath string    // Path to tmux server socket
	SocketName string    // Name of created tmux socket
	Sessions   []Session // List of sessions used on server initialization
}

// Creates a new server object.
func NewServer(socketPath, socketName string, sessions []Session) *Server {
	return &Server{
		SocketPath: socketPath,
		SocketName: socketName,
		Sessions:   sessions,
	}
}

// Lists all sessions managed by this server.
func (s *Server) ListSessions() ([]Session, error) {
	args := []string{
		"list-sessions",
		"-F", "#{session_id}:#{session_name}"}
	if s.SocketPath != "" {
		args = append([]string{"-S", s.SocketPath}, args...)
	}
	if s.SocketName != "" {
		args = append([]string{"-L", s.SocketName}, args...)
	}

	out, _, err := RunCmd(args)
	if err != nil {
		return nil, err
	}

	outLines := strings.Split(out, "\n")
	sessions := []Session{}
	re := regexp.MustCompile(`\$([0-9]+):(.+)`)
	for _, line := range outLines {
		result := re.FindStringSubmatch(line)
		if len(result) < 3 {
			continue
		}
		id, err_atoi := strconv.Atoi(result[1])
		if err_atoi != nil {
			return nil, err_atoi
		}

		sessions = append(sessions, Session{Name: result[2], Id: id})
	}

	return sessions, nil
}

// Add session to server configuration. This will change only in-library
// server representation. Used for initial configuration before creating new
// server.
func (s *Server) AddSession(session Session) {
	s.Sessions = append(s.Sessions, session)
}

// Create new session with given name on this server.
//
// Session always will be detached after creation. Call AttachSession to attach
// it. If session already exists, this function return an error. Check session
// with HaveSession before running it if you need it.
func (s *Server) NewSession(name string, args ...string) (session Session, err error) {
	if checkSessionName(name) == false {
		return session, errors.New("Bad session name")
	}

	args = append([]string{
		"new-session",
		"-d",
		"-D",
		"-s", name,
		"-P",
		"-F", "#{session_id}:#{session_name}"},
		args...)

	out, err_out, err_exec := RunCmd(args)
	if err_exec != nil {
		// It's okay, if session already exists.
		if !strings.Contains(err_out, "exit status 1") {
			return session, err_exec
		}
	}

	re := regexp.MustCompile(`\$([0-9]+):(.+)`)
	result := re.FindStringSubmatch(out)
	if len(result) < 3 {
		return session, errors.New("Error creating session")
	}

	id, err_atoi := strconv.Atoi(result[1])
	if err_atoi != nil {
		return session, err_atoi
	}

	session = Session{Name: result[2], Id: id}

	// Append the Session struct to the Sessions slice in the Server struct.
	s.Sessions = append(s.Sessions, session)

	return session, nil
}

// Kills session with given name. If killed session not found, KillSession will
// not raise error and just do nothing.
func (s *Server) KillSession(name string) error {
	// Running "kill-session" without name causes killing attached session.
	if checkSessionName(name) == false {
		return errors.New("KillSession: Bad session name")
	}

	args := []string{"kill-session", "-t", name}

	if _, _, err := RunCmd(args); err != nil {
		return err
	}

	return nil
}

// Return true that session with given name is exsits on this server, false
// otherwise.
func (s *Server) HasSession(name string) (bool, error) {
	if checkSessionName(name) == false {
		return false, errors.New("Bad session name")
	}

	args := []string{"has-session", "-t", name}

	_, err_out, err := RunCmd(args)
	if strings.Contains(err_out, "can't find session") {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Return list with all panes managed by this server.
func (s *Server) ListPanes() ([]Pane, error) {
	return ListPanes([]string{"-a"})
}
