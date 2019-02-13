/*
Package launchctl wraps running the appropriate commands for running the "launchctl setenv" command for a specific
EnvVar "Name" and "Val", for setting the environ for non-shell apps.
*/
package launchctl

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	"github.com/spf13/viper"
)

// Command holds the configuration for running the command to run "launchctl setenv" for a specific "Name=Val" combination.
type Command struct {
	Env  environ.Environ
	Name string
	Val  string
}

// CommandLine returns the appropriate command-line for setting the "Name" and "Val" environ.
func (s *Command) CommandLine() []string {
	return []string{"launchctl", "setenv", s.Name, s.Val}
}

// Exec runs the "CommandLine()" command, and deals with errors.
// StdOut and StdErr for the command is sent to the user.
// If No-op mode has been set, the command isn't run, but what would have been run is logged.
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
