/*
Package getenv wraps running the appropriate commands for getting the user's environ.
*/
package getenv

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const marker = "=====###====="

// Command holds the configuration for running the command necessary to get the EnvVars from the user's shell.
type Command struct {
	Shell   string
	Env     environ.Environ
	command []string
}

// CommandLine returns the appropriate command-line to extract the environ from the user's shell.
// Currently supports "bash" and "zsh".
func (s *Command) CommandLine() ([]string, error) {
	switch {
	case strings.HasSuffix(s.Shell, "bash"):
		return []string{
			s.Shell, "--login", "-c", "echo '" + marker + "'; env",
		}, nil
	case strings.HasSuffix(s.Shell, "zsh"):
		return []string{
			s.Shell, "-i", "--login", "-c", "echo '" + marker + "'; env",
		}, nil
	}
	log.WithField("Shell", s.Shell).Debug("Unrecognised shell")
	return nil, errors.Errorf("I don't know how to work with %#v", s.Shell)
}

// Exec runs the necessary command, with the given environ, and returns the Stdout output.
// StdErr for the command is sent to the user.
// If No-op mode has been set, the command isn't run, but what would have been run is logged.
func (s *Command) Exec() ([]byte, error) {
	command, err := s.CommandLine()
	if err != nil {
		return nil, err
	}
	if viper.GetBool("noop") {
		log.WithField("Command", command).Info("NOOP: Not running command")
		return nil, nil
	}

	log.WithField("Command", command).Info("Running command...")
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = s.Env.ToOsEnviron()
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		log.WithFields(log.Fields{"Command": command, "Cmd": cmd}).Debug("Command failed")
		return nil, errwrap.Wrap(errors.Errorf("%v failed to run", command), err)
	}

	log.WithFields(log.Fields{"Command": command, "Output": string(output[:])}).Debug("Command successful")
	return output, nil
}

// ExecToLines runs the command, and parses Stdout into an array of NL-separated lines.
func (s *Command) ExecToLines() ([]string, error) {
	output, err := s.Exec()
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	skipUntilMarker := true
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == marker:
			// This will mean the following lines will be collected, but not the current 'marker' line:
			skipUntilMarker = false
		case !skipUntilMarker:
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.WithFields(log.Fields{"Err": err, "Lines (so far)": lines}).Debug("Failed to parse the output")
		return nil, errwrap.Wrap(errors.Errorf("Failed to parse the command output into lines"), err)
	}

	log.WithField("Lines", lines).Debug("Parsed into lines")
	return lines, nil
}
