package main

import (
    "fmt"
    "net/http"
    "sync"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("RequestReader", func() {

    var list []string
    var reader *RequestReader

    BeforeEach(func() {
        list = []string{
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/1")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/2")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/3")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/4")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/5")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/6")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/7")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/8")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/9")),
            fmt.Sprintf(`%s -X POST `, UrlForTestServer("/10")),
        }
        file := CreateFileFromLines(list)
        file.Close()
        reader = NewRequestReader(file.Name())
    })

    It("Single reader iterates over lines in a file", func() {
        requests := []*http.Request{}
        stream := NewSequentialRequestStream(reader)
        for stream.HasNext() {
            req,err := stream.Next()
            if err != nil{
                panic(err)
            }
            requests = append(requests, req)
        }
        Expect(len(requests)).To(Equal(len(list)))
    })

    for _, numberOfWorkers := range NUMBER_OF_WORKERS_TO_TEST {
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
                        req,err := stream.Next()
                        if err != nil{
                            panic(err)
                        }
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
