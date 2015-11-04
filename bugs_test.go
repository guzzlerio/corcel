package main

import (
    "os"
    "path/filepath"
    "fmt"
    "strconv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bugs replication", func() {

	var (
		exePath string
		err     error
	)

	BeforeEach(func() {
		os.Remove("./output.yml")
		exePath, err = filepath.Abs("./corcel")
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		TestServer.Clear()
	})

    It("Error when running a simple run with POST and PUT #18", func() {
        numberOfWorkers := 2
        list := []string{
            fmt.Sprintf(`%s -X POST -d '{"name": "bob"}' -H "Content-type: application/json"`, UrlForTestServer("/success")),
            fmt.Sprintf(`%s -X PUT -d '{"id": 1,"name": "bob junior"}' -H "Content-type: application/json"`, UrlForTestServer("/success")),
            fmt.Sprintf(`%s?id=1 -X GET -H "Content-type: application/json"`, UrlForTestServer("/success")),
            fmt.Sprintf(`%s?id=1 -X DELETE -H "Content-type: application/json"`, UrlForTestServer("/success")),
        }

        SutExecute(list[:1], "--random", "--summary","--workers",strconv.Itoa(numberOfWorkers))

        var executionOutput ExecutionOutput
        UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(2)))
    })
})
