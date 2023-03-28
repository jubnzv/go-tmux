// The MIT License (MIT)
// Copyright (C) 2019-2023 Georgiy Komarov <jubnzv@gmail.com>

// TODO: TMUX_TMPDIR

package tmux

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
)

// Wrapper to tmux CLI that execute command with given arguments and returns
// stdout and stderr output.
func RunCmd(args []string) (string, string, error) {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return "", "", err
	}
	cmd := exec.Command(tmux, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	return outStr, errStr, err
}

// Execute tmux command using syscall execve(2).
func ExecCmd(args []string) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	args = append([]string{tmux}, args...)

	if err := syscall.Exec(tmux, args, os.Environ()); err != nil {
		return err
	}

	return nil
}

// Returns true if executed inside tmux, false otherwise.
func IsInsideTmux() bool {
	if os.Getenv("TMUX") != "" {
		return true
	} else {
		return false
	}
}
