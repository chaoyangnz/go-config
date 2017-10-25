/*
Package config abstracts configuring the config.logging and config.settings libraries.

You're expected to initalise this by calling the Init() function with a Config{}
struct.  The struct needs to have values set in it for configuring the above
libraries.  Alternatively, you can set the `FromConfig` setting, and it will
try to self-configure via the Viper script.
*/
package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/mexisme/go-config/logging"
	"github.com/mexisme/go-config/settings"
)

type Config struct {
	// File is the name of a file that Viper will read for configuration.
	// It searches for the file in the user's `$HOME` dir as well as the current working dir.
	File string
	// EnvPrefix is a required prefix-string that Viper uses to filter Env-vars
	// for settings.
	EnvPrefix string
	// Debug enables debug logging if set to `true`:
	Debug bool
	// FromConfig enables the following settings (Name ... LoggingSentryDsn)
	// to be configured via Viper.
	// This means it will use the above Config file and appropriate Env-vars
	FromConfig bool

	// Name is the App name, used in log messages
	Name string
	// Environment is the App's environment it was run in -- e.g. "staging" or "prod"
	Environment string
	// Release is the App's release / version string
	Release string
	// LoggingFormat sets the log-out format for log messages
	LoggingFormat string
	// LoggingSentryDsn is the connection string (DSN) used to send errors to Sentry.io
	LoggingSentryDsn string

	// We don't want to try to reinitialise the config more than once
	initConfigDone bool
	// We don't want to try to reinitialise the logging more than once
	logConfigDone bool
}

var config Config

// Init is to allow other packages to easily depend on this one,
// since most of the important logic is in init()
func Init(initConfig Config) {
	config = initConfig

	config.read()
	config.logging()
}

// DryRun says whether the dry_run config has been set
func DryRun(reason string, args ...interface{}) bool {
	dryRun := settings.DryRun()
	if dryRun {
		log.Infof("DRY-RUN MODE: "+reason, args...)
	}

	return dryRun
}

// AddConfigItems passes the configItems through to settings.AddConfigItems()
func AddConfigItems(configItems []string) {
	// Need to ensure the system has been configured at least once!
	config.read() // TODO: Viper dynamically reads -- this may not be needed.
	settings.AddConfigItems(configItems)
}

// ApplyWith passes the configItems through to settings.ApplyWith()
func ApplyWith(item string, f func(interface{})) {
	settings.ApplyWith(item, f)
}

func (s *Config) read() {
	// This should make it safe to rerun a few times
	if !s.initConfigDone {
		settings.ReadConfig(s.File, s.EnvPrefix)
		s.initConfigDone = true
	}
}

func (s *Config) logging() {
	s.read()

	// This should make it safe to rerun a few times
	if !s.logConfigDone {
		logConfig := logging.New()
		logConfig.SetAppName(s.Name).SetAppEnv(s.Environment).SetAppRelease(s.Release)
		logConfig.SetFormat(s.LoggingFormat).SetSentryDsn(s.LoggingSentryDsn)

		if s.FromConfig {
			logConfig.SetFromConfig()
		}

		logConfig.Init()

		s.logConfigDone = true
	}
}
