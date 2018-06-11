package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	"github.com/spf13/viper"
)

const (
	defaultShell = "zsh"
)

var (
	// TODO: Make this overridable:
	copiedEnvs = []string{"HOME", "LOGNAME", "USER", "LANG"}
)

func init() {
	settings()
	logging()
}

func settings() {
	viper.BindEnv("debug")
	viper.BindEnv("noop")
	// This allows us to override in a config file:
	viper.BindEnv("shell")

	// This means any "." chars in a FQ config name will be replaced with "_"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(".osx-env-sync")
	// TODO: Is this the default?
	viper.AddConfigPath("$HOME")

	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Debug("Using file")

	} else {
		log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Warn(err)
	}
}

func logging() {
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}
	if viper.GetBool("noop") {
		log.Debug("No-op mode enabled")
	}
}

/**
 * Enviroment Variables:
 */

type Environ map[string]string

func NewEnviron() Environ {
	return make(Environ)
}

func (s Environ) AddOsEnviron() Environ {
	return s.AddEnviron(os.Environ())
}

func (s Environ) AddEnviron(environ []string) Environ {
	for _, envLine := range environ {
		kv := strings.SplitN(envLine, "=", 2)
		s[kv[0]] = kv[1]
	}

	log.WithField("Env", s).Debug("Environ imported")
	return s
}

func (s Environ) NewFromFiltered(names []string) Environ {
	newEnv := make(Environ)

	for _, name := range names {
		newEnv[name] = s[name]
	}

	log.WithField("Filtered Env", newEnv).Debug("New filtered environ created")
	return newEnv
}

func (s Environ) ToOsEnviron() []string {
	newEnv := make([]string, 0)

	for name, val := range s {
		line := fmt.Sprintf("%s=%s", name, val)
		newEnv = append(newEnv, line)
	}

	log.WithField("os.Environ", newEnv).Debug("os.Environ created")
	return newEnv
}

/**
 * GetEnvCommand
 */

type Shell struct {
	Shell string
	Env   Environ
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

type GetEnvCommand struct {
	Shell   string
	Env     Environ
	command []string
}

func (s *GetEnvCommand) CommandLine() ([]string, error) {
	switch {
	case strings.HasSuffix(s.Shell, "bash"):
		return []string{s.Shell, "--login", "-c", "env"}, nil
	case strings.HasSuffix(s.Shell, "zsh"):
		return []string{s.Shell, "-i", "--login", "-c", "env"}, nil
	}
	log.WithField("Shell", s.Shell).Debug("Unrecognised shell")
	return nil, errors.Errorf("I don't know how to work with %#v", s.Shell)
}

func (s *GetEnvCommand) Exec() ([]byte, error) {
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

func (s *GetEnvCommand) ExecToLines() ([]string, error) {
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

type LaunchctlCommand struct {
	Env  Environ
	Name string
	Val  string
}

func (s *LaunchctlCommand) CommandLine() []string {
	return []string{"launchctl", "setenv", s.Name, s.Val}
}

func (s *LaunchctlCommand) Exec() error {
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

/**
 * Main
 */

func main() {
	env := NewEnviron().AddOsEnviron()
	log.WithField("Env", env).Debug()

	shell := (&Shell{Env: env}).Choose()
	log.WithField("Found Shell", shell).Debug()

	filteredEnv := env.NewFromFiltered(copiedEnvs)
	filteredEnv["TERM"] = "xterm"
	log.WithField("Filtered Env", filteredEnv).Debug()

	cmd := &GetEnvCommand{Shell: shell, Env: filteredEnv}

	envLines, err := cmd.ExecToLines()
	if err != nil {
		// fmt.Println(err.(*errors.Error).ErrorStack())
		log.Fatal(err)
	}

	givenEnv := NewEnviron().AddEnviron(envLines)
	log.WithField("Env", givenEnv).Debug("From the Shell")

	for k, v := range givenEnv {
		cmd := &LaunchctlCommand{Env: env, Name: k, Val: v}
		if err := cmd.Exec(); err != nil {
			log.Fatal(err)
		}
	}
}
