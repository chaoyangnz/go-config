package config_test

import (
	. "github.com/mexisme/go-config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config", func() {
	Describe("from a config file", func() {
		It("is set if AddConfigItems() happens before reading", func() {
			valInternal := ""

			AddConfigItems([]string{"test-before.test_str"})

			Init(Config{
				File:       "config-before",
				Dir:        "fixtures",
				FromConfig: true,
			})

			ApplyWith("test-before.test_str", func(val interface{}) {
				valInternal = val.(string)
			})

			Expect(valInternal).To(Equal("wootless"))
		})
	})
})
