package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/mexisme/go-config/logging"
	"github.com/mexisme/go-config/settings"
)

type Config struct {
	File string
	EnvPrefix string
	Name string
	Environment string
	Release string
	LoggingFormat string
	LoggingSentryDsn string

	Debug bool
	FromConfig bool

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
