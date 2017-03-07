package utils_test

import (
	"time"

	. "github.com/guzzlerio/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeUtils", func() {
	It("Time a function", func() {
		var result = Time(func() {})

		Expect(result).To(BeNumerically(">", time.Duration(1)))
	})

	It("DurationIsBetween succeeds", func() {
		var a = time.Duration(1)
		var b = time.Duration(2)
		var c = time.Duration(3)

		Expect(DurationIsBetween(b, a, c)).To(BeTrue())
	})

	It("DurationIsBetween fails", func() {
		var a = time.Duration(1)
		var b = time.Duration(2)
		var c = time.Duration(3)

		Expect(DurationIsBetween(a, b, c)).To(BeFalse())
	})
})
