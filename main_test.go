package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	server *RequestRecordingServer
	port   int
)

func TestMain(m *testing.M) {
	port = 8000
	server = CreateRequestRecordingServer(port)
	server.Start()
	os.Exit(m.Run())
	server.Stop()
}

var _ = Describe("Main", func() {

	BeforeEach(func() {
		server.Clear()
	})

	It("Makes a http POST request", func() {
		list := `http://127.0.0.1:8000/A -X POST`

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
		output, err := cmd.CombinedOutput()

		fmt.Println(string(output))

		Expect(err).To(BeNil())
		Expect(server.Contains(RequestWithPath("/A"), RequestWithMethod("POST"))).To(Equal(true))
	})

	It("Makes a http get request to each url in a file", func() {
		list := `http://127.0.0.1:8000/A
			http://127.0.0.1:8000/B
			http://127.0.0.1:8000/C`
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

		Expect(err).To(BeNil())
		Expect(server.Contains(RequestWithPath("/A"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/B"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/C"))).To(Equal(true))
	})
})
