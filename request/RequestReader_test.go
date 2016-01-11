package request

import (
	"fmt"
	"net/http"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "ci.guzzler.io/guzzler/corcel/utils"
)

var _ = Describe("RequestReader", func() {

	var list []string
	var reader *RequestReader

	BeforeEach(func() {
		list = []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/2")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/3")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/4")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/5")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/6")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/7")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/8")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/9")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/10")),
		}
		file := CreateFileFromLines(list)
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing the file")
		}
		reader = NewRequestReader(file.Name())
	})

	It("Single reader iterates over lines in a file", func() {
		requests := []*http.Request{}
		stream := NewSequentialRequestStream(reader)
		for stream.HasNext() {
			req, err := stream.Next()
			check(err)
			requests = append(requests, req)
		}
		Expect(len(requests)).To(Equal(len(list)))
	})

	for _, numberOfWorkers := range NumberOfWorkersToTest {
		It(fmt.Sprintf("Multiple readers totalling %v iterates over lines in a file", numberOfWorkers), func() {
			var wg sync.WaitGroup
			var mutex = &sync.Mutex{}
			wg.Add(numberOfWorkers)
			requests := []*http.Request{}
			for i := 0; i < numberOfWorkers; i++ {
				go func() {
					stream := NewSequentialRequestStream(reader)
					for stream.HasNext() {
						mutex.Lock()
						req, err := stream.Next()
						check(err)
						requests = append(requests, req)
						mutex.Unlock()
					}
					wg.Done()
				}()
			}
			wg.Wait()
			Expect(len(requests)).To(Equal(len(list) * numberOfWorkers))
		})
	}
})
