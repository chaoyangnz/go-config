package config_test

import (
	. "github.com/mexisme/go-config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config", func() {
	It("compiles", func() {
		Init(Config{})
	})

	Describe(".FromStringOrFunc", func() {
		It("returns the string value, when provided a string", func() {
			res, err := FromStringOrFunc("Hello!")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("Hello!"))
		})

		It("returns the func's return value, when provided a func", func() {
			res, err := FromStringOrFunc(func() string { return "Hello func!" }())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("Hello func!"))
		})

		It("fail if the type is not string or func()", func() {
			_, err := FromStringOrFunc(2000)
			Expect(err).To(HaveOccurred())
		})
	})
})
