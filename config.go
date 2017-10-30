/*
Package "config" provides an abstraction of a common pattern I use for configuring my application.
This includes configuring logging (via the "Logrus" package) so I can log debug-level information
around configuration at the same time.

The "config" package is a simple abstraction around the "settings" and "logging" sub-packages.

Initialisation

You're expected to initalise this by calling the Init() function with a Config{}
struct (defined below).  The struct needs to have values set in it for configuring the above
libraries.  Alternatively, you can enable the "FromConfig" setting, and it will
try to self-configure via the Viper script.
*/
package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/mexisme/go-config/logging"
	"github.com/mexisme/go-config/settings"
)

/*
Config provides basic fields for configuring the "settings" and "logging" packages.

"File" is the name of a file that Viper will read for configuration.
By default, it searches for the file in the user's `$HOME` dir as well as the current
workig dir -- but see "OnlyUseDir" below.
If the file-name is empty, settings won't be loaded from a file (only env-vars).

"Dir" is an optional additional additional dir to search for a config file.

"OnlyUseDir" when false will additionally search "$HOME" and current working dir for
the config file. When true, will only search in the above "Dir" directory.
If "Dir" is not given, then the config file won't be loaded.

"EnvPrefix" is a required prefix-string that Viper uses to filter Env-vars
for settings.

"Debug" enables debug logging if set to "true":

"FromConfig" enables the following settings (Name ... LoggingSentryDsn)
to be configured via Viper.
This means it will use the above Config file and appropriate Env-vars

"Name" is the App name, used in log messages.
This can be a string, or a func ref that will return a string.

"Environment" is the App's environment it was run in -- e.g. "staging" or "prod"
This can be a string, or a func ref that will return a string.

"Release" is the App's release / version string
This can be a string, or a func ref that will return a string.

"LoggingFormat" sets the log-out format for log messages

"LoggingSentryDsn" is the connection string (DSN) used to send errors to Sentry.io
*/
type Config struct {
	File string
	Dir string
	OnlyUseDir bool
	EnvPrefix string
	Debug bool

	FromConfig bool

	Name interface{}
	Environment interface{}
	Release interface{}
	LoggingFormat string
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
	settings.AddConfigItems(configItems)
}

// ApplyWith passes the configItems through to settings.ApplyWith()
func ApplyWith(item string, f func(interface{})) {
	settings.ApplyWith(item, f)
}

func (s *Config) read() {
	// This should make it safe to rerun a few times
	if !s.initConfigDone {
		settings.ReadConfig(s.File, s.Dir, s.EnvPrefix, s.OnlyUseDir)
		s.initConfigDone = true
	}
}

// FromStringOrFunc will return a different value depending on the provided val:
// - If it's a string, provide the given val
// - If it's a func(), provide teh val returned by the func
func FromStringOrFunc(val interface{}) (string, error) {
	switch val.(type) {
	case string:
		return val.(string), nil
	case func() string:
		f := val.(func() string)
		return f(), nil
	}

	return "", fmt.Errorf("Can't read value from %#v", val)
}

func (s *Config) logging() {
	s.read()

	// This should make it safe to rerun a few times
	if !s.logConfigDone {
		logConfig := logging.New()

		name, _ := FromStringOrFunc(s.Name)
		env, _ := FromStringOrFunc(s.Environment)
		release, _ := FromStringOrFunc(s.Release)

		logConfig.SetAppName(name).SetAppEnv(env).SetAppRelease(release)
		logConfig.SetFormat(s.LoggingFormat).SetSentryDsn(s.LoggingSentryDsn)

		if s.FromConfig {
			logConfig.SetFromConfig()
		}

		logConfig.Init()

		s.logConfigDone = true
	}
}
