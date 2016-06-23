package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	. "ci.guzzler.io/guzzler/corcel"
	"ci.guzzler.io/guzzler/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanHttpRequest", func() {
	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	It("Supplies a payload to the HTTP Request", func() {
		planBuilder := test.NewYamlPlanBuilder()

		path := "/people"
		body := "Zee Body"

		planBuilder.
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPRequestAction().URL(TestServer.CreateURL(path)).Body(body).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(body))).To(Equal(true))
	})

	It("Supplies a payload as a file reference to the HTTP Request", func() {
		content := []byte("temporary file's content")
		dir, err := ioutil.TempDir("", "ExecutionPlanHttpRequest")
		if err != nil {
			panic(err)
		}

		defer func() {
			err := os.RemoveAll(dir) // clean up
			if err != nil {
				panic(err)
			}
		}()

		tmpfn := filepath.Join(dir, "tmpfile")
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			panic(err)
		}

		path := "/people"
		body := fmt.Sprintf("@%s", tmpfn)

		planBuilder := test.NewYamlPlanBuilder()

		planBuilder.
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPRequestAction().URL(TestServer.CreateURL(path)).Body(body).Build())

		err = ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(string(content)))).To(Equal(true))
	})

	PIt("Returns an error when the file path does not exist", func() {})
})
