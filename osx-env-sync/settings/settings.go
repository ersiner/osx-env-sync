/*
Package settings is for setting-up some basic settings, incl. standard Viper settings, and basic logging:

- debug ($DEBUG) -- enable debug mode

- noop ($NOOP) -- enable No-op mode (no sub-commands will actually be run)

- shell (no env-var) -- allows overriding the standard shell-discovery in the "shell.Choose()" function

These can all optionally be set in a "~/.osx-env-sync.toml" config file, or using the above run-time EnvVars.
*/
package settings

import (
	"github.com/mexisme/osx-env-sync/osx-env-sync/version"

	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	Setup()
}

// Setup configures standard Viper and Logrus settings
func Setup() {
	config()
	logging()
}

func config() {
	viper.BindEnv("debug")
	viper.BindEnv("noop")
	// This allows us to override in a config file:
	viper.SetDefault("shell", "")

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
	// TODO: Should this be a Debug message?
	log.Infof("## %#v release %v ##", version.Application(), version.Release())

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}
	if viper.GetBool("noop") {
		log.Debug("No-op mode enabled")
	}
}
