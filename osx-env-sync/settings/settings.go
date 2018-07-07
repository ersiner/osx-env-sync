package settings

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

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
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}
	if viper.GetBool("noop") {
		log.Debug("No-op mode enabled")
	}
}
