package main

import (
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/processor"
	. "ci.guzzler.io/guzzler/corcel/utils"
)

var _ = Describe("Bugs replication", func() {

	BeforeEach(func() {
		err := os.Remove("./output.yml")
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	It("Error when running a simple run with POST and PUT #18", func() {
		numberOfWorkers := 2
		list := []string{
			fmt.Sprintf(`%s -X POST -d '{"name": "bob"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X PUT -d '{"id": 1,"name": "bob junior"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s?id=1 -X GET -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s?id=1 -X DELETE -H "Content-type: application/json"`, URLForTestServer("/success")),
		}

		SutExecute(list[:1], "--random", "--summary", "--workers", strconv.Itoa(numberOfWorkers))

		var executionOutput processor.ExecutionOutput
		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(2)))
	})

	PIt("Error when too many workers specified causing too many open files #23", func() {
		numberOfWorkers := 100000000000
		list := []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
		}

		output, err := InvokeCorcel(list, "--workers", strconv.Itoa(numberOfWorkers), "--progress", "none")

		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring("Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers"))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		output, err := InvokeCorcel(list, "--progress", "none")

		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring("Your urls in the test specification must be valid urls"))
	})
})
