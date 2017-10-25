package config_test

import (
	. "github.com/mexisme/go-config"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("config", func() {
	It("compiles", func() {
		Init(Config{})
	})
})
