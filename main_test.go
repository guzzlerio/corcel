package main_test

import (
	. "github.com/reaandrew/code-named-something"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var _ = Describe("Main", func() {
	It("Makes a http get request to each url in a file", func() {
		port := 8000
		list := `http://127.0.0.1:8000/A
http://127.0.0.1:8000/B
http://127.0.0.1:8000/C
		`
		server := CreateRequestRecordingServer(port)
		defer server.Stop()
		server.Start()

		file, err := ioutil.TempFile(os.TempDir(), "prefix")
		if err != nil {
			panic(err)
		}
		file.WriteString(list)
		defer os.Remove(file.Name())

		exePath, err := filepath.Abs("./code-named-something")
		if err != nil {
			panic(err)
		}
		cmd := exec.Command(exePath, "-f", file.Name())
		_, _ = cmd.CombinedOutput()

		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/A", nil)
		_, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		Expect(err).To(BeNil())
		Expect(server.Contains(RequestWithPath("/A"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/B"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/C"))).To(Equal(true))
	})
})
