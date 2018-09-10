package config

import (
	"time"

	"github.com/mexisme/go-config/settings"
)

// AddConfigItems passes the configItems through to settings.AddConfigItems()
func AddConfigItems(configKeys []string) {
	settings.AddConfigItems(configKeys)
}

// ApplyWith passes the configItems through to settings.ApplyWith()
func ApplyWith(key string, f func(interface{})) {
	settings.ApplyWith(key, f)
}

// Get forwards to settings.Get
func Get(key string) interface{} {
	return settings.Get(key)
}

// GetBool forwards to settings.GetBool
func GetBool(key string) bool {
	return settings.GetBool(key)
}

// GetFloat64 forwards to settings.GetFloat64
func GetFloat64(key string) float64 {
	return settings.GetFloat64(key)
}

// GetInt forwards to settings.GetInt
func GetInt(key string) int {
	return settings.GetInt(key)
}

// GetString forwards to settings.GetString
func GetString(key string) string {
	return settings.GetString(key)
}

// GetStringMap forwards to settings.GetStringMap
func GetStringMap(key string) map[string]interface{} {
	return settings.GetStringMap(key)
}

// GetStringMapString forwards to settings.GetStringMapString
func GetStringMapString(key string) map[string]string {
	return settings.GetStringMapString(key)
}

// GetStringSlice forwards to settings.GetStringSlice
func GetStringSlice(key string) []string {
	return settings.GetStringSlice(key)
}

// GetTime forwards to settings.GetTime
func GetTime(key string) time.Time {
	return settings.GetTime(key)
}

// GetDuration forwards to settings.GetDuration
func GetDuration(key string) time.Duration {
	return settings.GetDuration(key)
}

// IsSet forwards to settings.IsSet
func IsSet(key string) bool {
	return settings.IsSet(key)
}

// AllSettings forwards to settings.AllSettings
func AllSettings() map[string]interface{} {
	return settings.AllSettings()
}
