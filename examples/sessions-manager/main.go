// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package main

import (
	// "errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	tmux "github.com/jubnzv/go-tmux"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

const (
	version = "1.0.0"
)

var (
	app     = kingpin.New("tmux-session", "A utility to manage tmux sessions.")
	verbose = app.Flag("verbose", "Enable additional output").Bool()

	save     = app.Command("save", "Write current tmux session in yaml")
	savePath = save.Arg("path", "Path to saving yaml configuration").String()
	fileName = save.Arg("name", "yaml file name").String()

	load     = app.Command("load", "Load tmux configuration from a yaml file")
	loadPath = load.Arg("path", "Path to yaml configuration").String()
)

// Represents a tmux window configuration
type WindowConf struct {
	WindowName     string `yaml:"window_name"`
	StartDirectory string `yaml:"start_directory"`
}

// Represents a tmux session configuration
type SessionConf struct {
	SessionName string       `yaml:"session_name"`
	Windows     []WindowConf `yaml:"windows"`
}

// Expand ~ in path.
func expandTilde(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, nil
}

// Generate tmux session name from yaml configuration filepath
func getSessionNameFromFilepath(fpath string) string {
	session_name := filepath.Base(fpath)
	session_name = session_name[0 : len(fpath)-len(filepath.Ext(session_name))]
	return session_name
}

// Saves current tmux session in yaml file
func doSave(savePath string, fileName string, verbose bool) error {
	// if !tmux.IsInsideTmux() {
	//     return errors.New("Not inside tmux.")
	// }

	// Get current session configuration
	session_name, err := tmux.GetAttachedSessionName()
	if err != nil {
		return err
	}
	session := tmux.Session{Name: session_name}
	windows, err := session.ListWindows()
	if err != nil {
		return err
	}

	// Generate filename if not specified
	if len(fileName) == 0 {
		fileName = session_name
	}

	// Collect windows configurations
	windows_conf := []WindowConf{}
	for _, w := range windows {
		windows_conf = append(windows_conf, WindowConf{
			WindowName:     w.Name,
			StartDirectory: w.StartDirectory})
		fmt.Println(w.StartDirectory)
	}

	// Prepare session configuration
	session_conf := SessionConf{
		SessionName: fileName,
		Windows:     windows_conf}

	conf_out, err := yaml.Marshal(&session_conf)
	if err != nil {
		return err
	}

	savePath, err = expandTilde(savePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, conf_out, 0644)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Saved in %s", savePath)
	}

	return nil
}

// Loads tmux session from json configuration file and attach to it.
func doLoad(loadPath string) error {
	// Parse configuration
	file_content, err := ioutil.ReadFile(loadPath)
	if err != nil {
		return err
	}
	conf_in := SessionConf{}
	err = yaml.Unmarshal(file_content, &conf_in)
	if err != nil {
		return err
	}

	// Prepare go-tmux session
	session := tmux.Session{Name: conf_in.SessionName}
	for idx, w := range conf_in.Windows {
		session.AddWindow(tmux.Window{
			Name:           w.WindowName,
            // TODO: Read base-index from .tmux.conf
			Id:             idx + 1,
			SessionId:      session.Id,
			SessionName:    session.Name,
			StartDirectory: w.StartDirectory,
		})
	}

	// Create new session
	server := new(tmux.Server)
	server.AddSession(session)
	sessions := []*tmux.Session{}
	sessions = append(sessions, &session)
	conf := tmux.Configuration{
		Server:        server,
		Sessions:      sessions,
		ActiveSession: nil}

	// Apply created session configuration
	err = conf.Apply()
	if err != nil {
		return err
	}

	// Attach to created session
	err = session.AttachSession()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	rc := 0
	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case save.FullCommand():
		// Prepare filepath
		var path_to_save string
		if *savePath == "" {
			session_name, err := tmux.GetAttachedSessionName()
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			path_to_save = fmt.Sprintf("~/.tmux/sessions/%s.yaml", session_name)
		} else {
			path_to_save = *savePath
		}

		// Perform saving
		err = doSave(path_to_save, *fileName, *verbose)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

	case load.FullCommand():
		// Perform loading
		err = doLoad(*loadPath)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	}

	if err != nil {
		rc = -1
	}

	os.Exit(rc)
}
