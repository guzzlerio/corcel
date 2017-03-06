package core_test

import (
	. "github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Registry", func() {

	It("Creates a Registry", func() {
		var registry = CreateRegistry()

		Expect(len(registry.AssertionParsers)).To(Equal(0))
		Expect(len(registry.ActionParsers)).To(Equal(0))
		Expect(len(registry.ResultProcessors)).To(Equal(0))
		Expect(len(registry.ExtractorParsers)).To(Equal(0))
	})

	It("AddExtractorParser", func() {
		var registry = CreateRegistry().AddExtractorParser(yaml.RegexExtractorParser{})
		Expect(len(registry.AssertionParsers)).To(Equal(0))
		Expect(len(registry.ActionParsers)).To(Equal(0))
		Expect(len(registry.ResultProcessors)).To(Equal(0))
		Expect(len(registry.ExtractorParsers)).To(Equal(1))
	})

	It("AddAssertionParser", func() {
		var registry = CreateRegistry().AddAssertionParser(yaml.ExactAssertionParser{})
		Expect(len(registry.AssertionParsers)).To(Equal(1))
		Expect(len(registry.ActionParsers)).To(Equal(0))
		Expect(len(registry.ResultProcessors)).To(Equal(0))
		Expect(len(registry.ExtractorParsers)).To(Equal(0))
	})

	It("AddActionParser", func() {
		var registry = CreateRegistry().AddActionParser(http.YamlHTTPRequestParser{})
		Expect(len(registry.AssertionParsers)).To(Equal(0))
		Expect(len(registry.ActionParsers)).To(Equal(1))
		Expect(len(registry.ResultProcessors)).To(Equal(0))
		Expect(len(registry.ExtractorParsers)).To(Equal(0))
	})
})
