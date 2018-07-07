package command

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"

	log "github.com/Sirupsen/logrus"
)

const (
	defaultShell = "zsh"
)

type Shell struct {
	Shell string
	Env   environ.Environ
}

func (s *Shell) Choose() string {
	if s.Shell != "" {
		log.WithField("Shell", s.Shell).Debug("Using provided shell")
		return s.Shell
	}
	if shell, ok := s.Env["SHELL"]; ok && shell != "" {
		log.WithField("Shell", shell).Debug("Using env-var $SHELL")
		return shell
	}
	log.WithField("Shell", defaultShell).Debug("Using default shell")
	return defaultShell
}
