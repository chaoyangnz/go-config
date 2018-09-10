package config_test

import (
	. "github.com/mexisme/go-config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config", func() {
	Describe("from a config file", func() {
		It("is set if AddConfigItems() happens after reading", func() {
			valInternal := ""

			Init(Config{
				File:       "config-after",
				Dir:        "fixtures",
				FromConfig: true,
			})

			AddConfigItems([]string{"test-after.test_str"})

			ApplyWith("test-after.test_str", func(val interface{}) {
				valInternal = val.(string)
			})

			Expect(valInternal).To(Equal("wootles"))
		})
	})
})
