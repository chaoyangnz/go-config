package settings

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"os"
	"strings"
)

var initConfigDone = false

// ReadConfig uses Viper to read the configuration from .config.* files or Env Vars
// TODO:  list config items
func ReadConfig() {
	viper.BindEnv("debug")
	viper.BindEnv("base")

	viper.BindEnv("application.name", "APPLICATION_NAME")
	viper.BindEnv("application.environment", "ENVIRONMENT")

	// This means any "." chars in a FQ config name will be replaced with "_"
	// e.g. "sentry.dsn" --> "$CONFIG_SENTRY_DSN" instead of "$CONFIG_SENTRY.DSN"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Debug("Using file")

	} else {
		log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Error(err)
	}
}

// AddConfigItems adds a new configuration item, and makes it overridable by env vars
func AddConfigItems(configItems []string) {
	for _, item := range configItems {
		viper.BindEnv(item)
	}
}
