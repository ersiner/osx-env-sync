package getenv

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	"bufio"
	"bytes"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	"github.com/spf13/viper"
)

type Command struct {
	Shell   string
	Env     environ.Environ
	command []string
}

func (s *Command) CommandLine() ([]string, error) {
	switch {
	case strings.HasSuffix(s.Shell, "bash"):
		return []string{s.Shell, "--login", "-c", "env"}, nil
	case strings.HasSuffix(s.Shell, "zsh"):
		return []string{s.Shell, "-i", "--login", "-c", "env"}, nil
	}
	log.WithField("Shell", s.Shell).Debug("Unrecognised shell")
	return nil, errors.Errorf("I don't know how to work with %#v", s.Shell)
}

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
	output, err := cmd.Output()
	if err != nil {
		log.WithFields(log.Fields{"Command": command, "Cmd": cmd}).Debug("Command failed")
		return nil, errwrap.Wrap(errors.Errorf("%v failed to run", command), err)
	}

	log.WithFields(log.Fields{"Command": command, "Output": string(output[:])}).Debug("Command successful")
	return output, nil
}

func (s *Command) ExecToLines() ([]string, error) {
	output, err := s.Exec()
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.WithFields(log.Fields{"Err": err, "Lines (so far)": lines}).Debug("Failed to parse the output")
		return nil, errwrap.Wrap(errors.Errorf("Failed to parse the command output into lines"), err)
	}

	log.WithField("Lines", lines).Debug("Parsed into lines")
	return lines, nil
}
