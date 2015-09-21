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

func CreateList(lines string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	file.WriteString(lines)
	return file
}

var _ = Describe("Main", func() {

	var (
		exePath string
		err     error
	)

	BeforeEach(func() {
		exePath, err = filepath.Abs("./code-named-something")
		if err != nil {
			panic(err)
		}
		server.Clear()
	})

	supportedMethods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range supportedMethods {
		It(fmt.Sprintf("Makes a http %s request", method), func() {
			list := fmt.Sprintf(`http://127.0.0.1:8000/A -X %s`, method)
			file := CreateList(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())
			Expect(server.Contains(RequestWithPath("/A"), RequestWithMethod(method))).To(Equal(true))
		})

	}

	It("Makes a http get request to each url in a file", func() {
		list := `http://127.0.0.1:8000/A
			http://127.0.0.1:8000/B
			http://127.0.0.1:8000/C`
		file := CreateList(list)
		defer os.Remove(file.Name())

		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		Expect(err).To(BeNil())
		Expect(server.Contains(RequestWithPath("/A"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/B"))).To(Equal(true))
		Expect(server.Contains(RequestWithPath("/C"))).To(Equal(true))
	})
})
