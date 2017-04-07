package utils_test

import (
	"errors"
	"testing"

	. "github.com/guzzlerio/corcel/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRuntimeUtils(t *testing.T) {
	Convey("RuntimeUtils", t, func() {
		//TODO: Understand how to make this work in GoConvey
		SkipConvey("CheckErr", func() {
			So(func() {
				CheckErr(errors.New("BANG"))
			}, ShouldPanicWith, errors.New("BANG"))
		})
	})
}
