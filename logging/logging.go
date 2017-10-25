package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/mexisme/go-config/settings"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configItemAppName = "application.name"
	configItemAppEnv = "application.environment"
	configItemFormat = "logging.format"
	configItemSentryDsn = "logging.sentry.dsn"
)

// Config contains the details for configuring Logging
type Config struct {
	appName string
	appEnv string
	appRelease string
	sentryDsn string
}

// Logging is a singleton for managing logging
var Logging *Config

// New creates a new struct for managing logging
func New() *Config {
	return &Config{}
}

func (s *Config) SetFromConfig() *Config {
	settings.AddConfigItems([]string{
		configItemAppName, configItemAppEnv,
		configItemFormat, configItemSentryDsn,
	})

	settings.ApplyWith(configItemAppName, func(val interface{}) {
		s.SetAppName(val.(string))
	})
	settings.ApplyWith(configItemAppEnv, func(val interface{}) {
		s.SetAppEnv(val.(string))
	})

	// We do this before setting debug mode, to help-out the log aggregators:
	settings.ApplyWith(configItemFormat, func(val interface{}) {
		s.SetFormat(val.(string))
	})
	settings.ApplyWith(settings.ConfigItemDebug, func(val interface{}) {
		s.SetDebug(val.(bool))
	})
	settings.ApplyWith(configItemSentryDsn, func(val interface{}) {
		s.SetSentryDsn(val.(string))
	})

	return s
}

// SetAppName TODO
func (s *Config) SetAppName(appName string) *Config {
	s.appName = appName
	return s
}

// SetAppEnv TODO
func (s *Config) SetAppEnv(appEnv string) *Config {
	s.appEnv = appEnv
	return s
}

// SetAppRelease TODO
func (s *Config) SetAppRelease(appRelease string) *Config {
	s.appRelease = appRelease
	return s
}

// SetDebug TODO
func (s *Config) SetDebug(debug bool) *Config {
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")

	} else {
		log.SetLevel(log.InfoLevel)
	}

	return s
}

// SetFormat TODO
func (s *Config) SetFormat(loggingFormat string) *Config {
	switch loggingFormat {
	case "":
		fallthrough

	case "text":
		// This is the default log formatter in logrus, anyway:
		log.SetFormatter(&log.TextFormatter{})

	case "json":
		log.SetFormatter(&log.JSONFormatter{})

	default:
		log.Panicf("Log format %#v not supported.", loggingFormat)
	}

	return s
}

// SetLoggingSentryDsn TODO
func (s *Config) SetSentryDsn(sentryDsn string) *Config {
	s.sentryDsn = sentryDsn
	return s
}

// Init set-ups the Logrus library -- debug mode, etc
// Currently set-up via Viper
func (s *Config) Init() {
	appName := s.appName
	appRelease := s.appRelease
	appEnv := s.appEnv

	if appName == "" {
		appName = os.Args[0]
	}
	if appRelease != "" {
		appRelease = fmt.Sprintf(" release %v", appRelease)
	}
	if appEnv != "" {
		appEnv = fmt.Sprintf(" (in %v)", appEnv)
	}

	// TODO: Should this be a Debug message?
	log.Infof("## %#v%v%v ##", appName, appRelease, appEnv)

	if s.sentryDsn != "" {
		if err := s.setupSentry(s.sentryDsn); err != nil {
			log.Error(err)
		}
	}
}

func (s *Config) setupSentry(sentryDsn string) error {
	log.WithFields(log.Fields{"sentry.dsn": sentryDsn}).Debug("Configuring connection to Sentry.io")

	// TODO: Meta-tag for environment
	// Some meta-tags
	tags := map[string]string{
		// TODO: Pick a better name, as this maps to "${CONFIG}_APP":
		"app": viper.GetString("app"),
	}

	// Sentry will only log for messages of the following severity:
	levels := []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	}

	// Hook Sentry into Logrus:
	hook, err := logrus_sentry.NewWithTagsSentryHook(sentryDsn, tags, levels)
	if err != nil {
		return err
	}

	if s.appRelease != "" {
		// Set the Sentry "release" version:
		log.WithFields(log.Fields{"release": s.appRelease}).Debug("Setting release version in Sentry")
		hook.SetRelease(s.appRelease)
	}
	if s.appEnv != "" {
		// Set the Sentry "environment":
		log.WithFields(log.Fields{"environment": s.appEnv}).Debug("Setting environment in Sentry")
		hook.SetEnvironment(s.appEnv)
	}

	hook.StacktraceConfiguration.Enable = true

	// It seems as if the default 100ms is too short:
	hook.Timeout = 1 * time.Second

	// Now, add it into the Logrus hook-chain
	log.AddHook(hook)

	log.Info("Sentry enabled")

	return nil
}
