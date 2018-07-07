package launchctl

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	"github.com/spf13/viper"
)

type Command struct {
	Env  environ.Environ
	Name string
	Val  string
}

func (s *Command) CommandLine() []string {
	return []string{"launchctl", "setenv", s.Name, s.Val}
}

func (s *Command) Exec() error {
	command := s.CommandLine()
	if viper.GetBool("noop") {
		log.WithField("Command", command).Info("NOOP: Not running command")
		return nil
	}

	log.WithField("Command", command).Info("Running command...")
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = s.Env.ToOsEnviron()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{"Command": command, "Cmd": cmd}).Debug("Command failed")
		return errwrap.Wrap(errors.Errorf("%v failed to run", command), err)
	}

	return nil
}
