/*
Package settings manages reading config from a config file or env vars.

Supported Settings

- debug ($DEBUG) -- enables debug logging.
- dry_run ($DRY_RUN) -- allows enabling of a "dry-run" mode

Note: other packages may add other settings.
*/
package settings

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// ConfigItemDebug defines the Viper config item for running in debug mode
	ConfigItemDebug = "debug"
	// ConfigItemDryRun defines the Viper config item for enabling dry-run mode
	ConfigItemDryRun  = "dry_run"
)

// ReadConfig uses Viper to read the configuration from .config.* files or Env Vars
// `configFile` is the default config file name
// `envPrefix` allows you to add a Viper "EnvPrefix" to config env-vars
// TODO:  list config items
func ReadConfig(configFile, envPrefix string) {
	// This means any "." chars in a FQ config name will be replaced with "_"
	// e.g. "sentry.dsn" --> "$SENTRY_DSN" instead of "$SENTRY.DSN" (which won't work)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.BindEnv(ConfigItemDebug)
	viper.BindEnv(ConfigItemDryRun)

	if configFile != "" {
		viper.SetConfigName(configFile)
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err == nil {
			log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Debug("Using file")

		} else {
			log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Error(err)
		}
	}
}

// DryRun says whether the dry_run config has been set
func DryRun() bool {
	// Note: Not being set should count as "false"
	return viper.GetBool(ConfigItemDryRun)
}

// AddConfigItems adds a new configuration item, and makes it overridable by env vars
func AddConfigItems(configItems []string) {
	for _, item := range configItems {
		viper.BindEnv(item)
	}
}

// ApplyWith gets a setting from viper, and passes it to a closure
func ApplyWith(item string, f func(interface{})) {
	if viper.IsSet(item) {
		f(viper.Get(item))
	}
}
