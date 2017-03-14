package utils_test

import (
	"net/http"

	. "github.com/guzzlerio/corcel/utils"
	"github.com/guzzlerio/rizo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpUtils", func() {
	It("ConcatRequestPaths", func() {

		req1, _ := http.NewRequest("GET", "http://a.com/A", nil)
		req2, _ := http.NewRequest("GET", "http://a.com/B", nil)
		req3, _ := http.NewRequest("GET", "http://a.com/C", nil)

		var result = ConcatRequestPaths([]*http.Request{req1, req2, req3})

		Expect(result).To(Equal("/A/B/C"))
	})

	It("ToHTTPRequestArray", func() {
		req1, _ := http.NewRequest("GET", "http://a.com/A", nil)
		req2, _ := http.NewRequest("GET", "http://a.com/B", nil)
		req3, _ := http.NewRequest("GET", "http://a.com/C", nil)

		var a = rizo.RecordedRequest{
			Request: req1,
		}
		var b = rizo.RecordedRequest{
			Request: req2,
		}
		var c = rizo.RecordedRequest{
			Request: req3,
		}

		var inArray = []rizo.RecordedRequest{a, b, c}
		var expectedArray = []*http.Request{req1, req2, req3}

		Expect(ToHTTPRequestArray(inArray)).To(Equal(expectedArray))
	})
})
