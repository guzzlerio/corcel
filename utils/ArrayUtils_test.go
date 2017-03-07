package utils_test

import (
	. "github.com/guzzlerio/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArrayUtils", func() {
	var input = []string{"A", "B", "C"}

	It("ContainsString suceeds", func() {
		Expect(ContainsString(input, "B")).To(BeTrue())
	})

	It("ContainsString fails", func() {
		Expect(ContainsString(input, "D")).To(BeFalse())
	})
})
