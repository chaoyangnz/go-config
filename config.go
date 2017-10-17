package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/mexisme/go-config/logging"
	"github.com/mexisme/go-config/settings"
)

var (
	// We don't want to try to reinitialise the config more than once
	initConfigDone = false
	// We don't want to try to reinitialise the logging more than once
	logConfigDone = false
)

func init() {
	readConfig()
	configLogging()
}

// ImportMe is to allow other packages to easily depend on this one,
// since most of the important logic is in init()
func ImportMe() {
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
	readConfig() // TODO: Viper dynamically reads -- this may not be needed.
	settings.AddConfigItems(configItems)
}

// ApplyWith passes the configItems through to settings.ApplyWith()
func ApplyWith(item string, f func(interface{})) {
	settings.ApplyWith(item, f)
}

func readConfig() {
	// This should make it safe to rerun a few times
	if !initConfigDone {
		settings.ReadConfig()
		initConfigDone = true
	}
}

func configLogging() {
	readConfig()

	// This should make it safe to rerun a few times
	if !logConfigDone {
		logging.Configure()
		logConfigDone = true
	}
}
