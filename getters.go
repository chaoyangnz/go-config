package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/mexisme/go-config/logging"
	"github.com/mexisme/go-config/settings"
)

// AddConfigItems passes the configItems through to settings.AddConfigItems()
func AddConfigItems(configItems []string) {
	settings.AddConfigItems(configItems)
}

// ApplyWith passes the configItems through to settings.ApplyWith()
func ApplyWith(item string, f func(interface{})) {
	settings.ApplyWith(item, f)
}
