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
	ConfigItemDryRun = "dry_run"
)

// ReadConfig uses Viper to read the configuration from .config.* files or Env Vars
// `configFile` is the default config file name
// `envPrefix` allows you to add a Viper "EnvPrefix" to config env-vars
// `useOnlyDir` disables looking for a config file in "$HOME" or "." directories.
// TODO:  list config items
func ReadConfig(configFile, configDir, envPrefix string, onlyUseDir bool, autoBindEnv bool) {
	// This means any "." chars in a FQ config name will be replaced with "_"
	// e.g. "sentry.dsn" --> "$SENTRY_DSN" instead of "$SENTRY.DSN" (which won't work)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.BindEnv(ConfigItemDebug)
	viper.BindEnv(ConfigItemDryRun)

	if configFile != "" && !(onlyUseDir && configDir == "") {
		viper.SetConfigName(configFile)

		// Set the config dir's to search for the given file-name.
		// It selects the first one it finds:
		if configDir != "" {
			viper.AddConfigPath(configDir)
		}
		if !onlyUseDir {
			viper.AddConfigPath("$HOME")
			viper.AddConfigPath(".")
		}

		if err := viper.ReadInConfig(); err == nil {
			log.WithFields(log.Fields{"config_file": viper.ConfigFileUsed()}).Debug("Using file")
			if autoBindEnv {
				for _, key := range viper.AllKeys() {
					viper.BindEnv(key)
				}
			}
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
