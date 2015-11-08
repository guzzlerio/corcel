package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
)

var _ = Describe("Bugs replication", func() {

	BeforeEach(func() {
		os.Remove("./output.yml")
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

		var executionOutput ExecutionOutput
		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(2)))
	})

	It("Error when too many workers specified causing too many open files #23", func() {
		numberOfWorkers := 100000000000
		list := []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
		}

		output, err := InvokeCorcel(list, "--workers", strconv.Itoa(numberOfWorkers))

		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring("Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers"))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		output, err := InvokeCorcel(list)

		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring("Your urls in the test specification must be valid urls"))
	})
})
