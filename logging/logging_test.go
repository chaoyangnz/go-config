package logging_test

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	. "github.com/mexisme/go-config/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config/logging", func() {
	It("compiles", func() {
		// Only here to avoid failing the gomega import:
		Expect(true).To(BeTrue())
	})

	Describe("configures", func() {
		var buf *bytes.Buffer
		BeforeEach(func() {
			buf = new(bytes.Buffer)
			log.SetOutput(buf)

			viper.Set("logging.format", "text")
		})

		It("with text output", func() {
			New().SetFromConfig().Init()

			Expect(buf.String()).To(MatchRegexp("logging\\.test"))
		})

		It("with correct app name and env", func() {
			viper.Set("application.name", "APPL")
			viper.Set("application.environment", "ENV")
			New().SetFromConfig().Init()

			Expect(buf.String()).To(MatchRegexp("msg=.*APPL.*(in ENV)"))
		})
	})

	Describe("debug setting", func() {
		var buf *bytes.Buffer
		BeforeEach(func() {
			buf = new(bytes.Buffer)
			log.SetOutput(buf)

			viper.Set("logging.format", "text")
		})

		It("is switched-off correctly from '0'", func() {
			viper.Set("debug", 0)
			New().SetFromConfig().Init()

			log.Debug("Out it comes!")
			Expect(buf.String()).NotTo(MatchRegexp("Out it comes!"))
		})

		It("is switched-on correctly from '1'", func() {
			viper.Set("debug", 1)
			New().SetFromConfig().Init()

			log.Debug("Out it comes!")
			Expect(buf.String()).To(MatchRegexp("Out it comes!"))
		})

		It("is switched-off correctly from ''", func() {
			viper.Set("debug", "")
			New().SetFromConfig().Init()

			log.Debug("Out it comes!")
			Expect(buf.String()).NotTo(MatchRegexp("Out it comes!"))
		})

		It("is switched-on correctly from 'Hello'", func() {
			viper.Set("debug", "Hello")
			New().SetFromConfig().Init()

			log.Debug("Out it comes!")
			Expect(buf.String()).To(MatchRegexp("Out it comes!"))
		})
	})
})
