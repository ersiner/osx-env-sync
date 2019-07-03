package main

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/command"
	"github.com/mexisme/osx-env-sync/osx-env-sync/command/getenv"
	"github.com/mexisme/osx-env-sync/osx-env-sync/command/launchctl"
	"github.com/mexisme/osx-env-sync/osx-env-sync/environ"
	_ "github.com/mexisme/osx-env-sync/osx-env-sync/settings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// TODO: Make this overridable:
	copiedEnvs = []string{"HOME", "LOGNAME", "USER", "LANG"}
)

func main() {
	env := environ.New().AddOsEnviron()
	log.WithField("Env", env).Debug()

	shell := (&command.Shell{Shell: viper.GetString("shell"), Env: env}).Choose()
	log.WithField("Found Shell", shell).Debug()

	filteredEnv := env.NewFromFiltered(copiedEnvs)
	filteredEnv["TERM"] = "xterm"
	log.WithField("Filtered Env", filteredEnv).Debug()

	cmd := &getenv.Command{Shell: shell, Env: filteredEnv}

	envLines, err := cmd.ExecToLines()
	if err != nil {
		// fmt.Println(err.(*errors.Error).ErrorStack())
		log.Fatal(err)
	}

	givenEnv := environ.New().AddEnviron(envLines)
	log.WithField("Env", givenEnv).Debug("From the Shell")

	for k, v := range givenEnv {
		cmd := &launchctl.Command{Env: env, Name: k, Val: v}
		if err := cmd.Exec(); err != nil {
			log.Fatal(err)
		}
	}
}
