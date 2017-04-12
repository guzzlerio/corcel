package http_test

import (
	nethttp "net/http"
	"testing"

	"github.com/guzzlerio/corcel/infrastructure/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpDefaultsExtractor(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Http DefaultsExtractor", t, func() {
		Convey("Sets Defaults", func() {
			var result http.HttpActionState
			func() {
				var input = map[string]interface{}{}
				var extractor = http.NewDefaultsExtractor()

				result = extractor.Extract(input)
			}()
			Convey("sets default URL to empty string", func() {
				So(result.URL, ShouldEqual, "")
			})
			Convey("sets default method to GET", func() {
				So(result.Method, ShouldEqual, "GET")
			})
			Convey("sets default body to empty string", func() {
				So(result.Body, ShouldEqual, "")
			})
			Convey("sets default header to empty collection", func() {
				So(result.Headers, ShouldResemble, nethttp.Header{})
			})
		})

		Convey("Extracts HttpActionState", func() {
			var input = map[string]interface{}{}
			input["defaults"] = map[string]interface{}{}
			var defaults = input["defaults"].(map[string]interface{})
			defaults["HttpAction"] = map[string]interface{}{}

			var action = defaults["HttpAction"].(map[string]interface{})
			action["headers"] = map[string]interface{}{}

			var headers = action["headers"].(map[string]interface{})

			headers["key"] = "value"

			action["method"] = "GET"
			action["body"] = "Bang Bang"
			action["url"] = "http://somewhere"

			var extractor = http.NewDefaultsExtractor()

			var result = extractor.Extract(input)

			So(result.Headers.Get("key"), ShouldEqual, "value")
			So(result.Method, ShouldEqual, action["method"])
			So(result.Body, ShouldEqual, action["body"])
			So(result.URL, ShouldEqual, action["url"])
		})

		Convey("returns Empty state when no defaults", func() {
			var extractor = http.NewDefaultsExtractor()
			var state = map[string]interface{}{}
			So(extractor.Extract(state), ShouldNotBeNil)
		})

		Convey("returns Empty state when no default HttpAction", func() {
			var extractor = http.NewDefaultsExtractor()
			var state = map[string]interface{}{}
			state["defaults"] = map[string]interface{}{}
			So(extractor.Extract(state), ShouldNotBeNil)
		})

		Convey("returns Empty state when no http definitions", func() {
			var extractor = http.NewDefaultsExtractor()
			var state = map[string]interface{}{}
			state["defaults"] = map[string]interface{}{}
			var defaults = state["defaults"].(map[string]interface{})
			defaults["HttpAction"] = map[string]interface{}{}
			So(extractor.Extract(state), ShouldNotBeNil)
		})
	})
}
