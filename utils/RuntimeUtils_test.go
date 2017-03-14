package utils_test

import (
	"errors"

	. "github.com/guzzlerio/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuntimeUtils", func() {
	It("CheckErr", func() {
		Expect(func() {
			CheckErr(errors.New("BANG"))
		}).To(Panic())
	})
})
