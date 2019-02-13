/*
Package command contains some generic methods for working with command-lines.
*/
package command

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultShell for when no other is provided.
	DefaultShell = "zsh"
)

// Shell type contains some settings for determining the preferred shell.
type Shell struct {
	Shell string
	Env   environ.Environ
}

/*
Choose selects which shell should be used, from the first non-null, non-empty of the following:
- The value in the "Shell" field
- The value in the "$SHELL" Env-var
- The value of the "DefaultShell"
*/
func (s *Shell) Choose() string {
	if s.Shell != "" {
		log.WithField("Shell", s.Shell).Debug("Using provided shell")
		return s.Shell
	}
	if shell, ok := s.Env["SHELL"]; ok && shell != "" {
		log.WithField("Shell", shell).Debug("Using env-var $SHELL")
		return shell
	}
	log.WithField("Shell", DefaultShell).Debug("Using default shell")
	return DefaultShell
}
