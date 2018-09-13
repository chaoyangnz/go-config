package config_test

import (
	. "github.com/mexisme/go-config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("config", func() {
	Describe("from a config file", func() {
		It("is not to bind environment variables automatically after reading and initialising with AutoBindEnv off", func() {
			valInternal := ""
			os.Setenv("TEST_FOO_TEST_BAR", "test-bar-value-from-env-var")

			Init(Config{
				File:       "config-auto-bind-env",
				Dir:        "fixtures",
				FromConfig: true,
			})

			ApplyWith("test_foo.test_bar", func(val interface{}) {
				valInternal = val.(string)
			})
			Expect(valInternal).To(Equal("test-bar-value"))
		})

		It("is to bind environment variables automatically after reading and initialising with AutoBindEnv on", func() {
			valInternal := ""
			os.Setenv("TEST_FOO_TEST_BAR", "test-bar-value-from-env-var")

			Init(Config{
				File:       "config-auto-bind-env",
				Dir:        "fixtures",
				FromConfig: true,
				AutoBindEnv: true,
			})

			ApplyWith("test_foo.test_bar", func(val interface{}) {
				valInternal = val.(string)
			})
			Expect(valInternal).To(Equal("test-bar-value-from-env-var"))
		})
	})
})
