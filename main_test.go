package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func containsString(list []string, expected string) bool {
	for _, current := range list {
		if current == expected {
			return true
		}
	}
	return false
}

var _ = Describe("Main", func() {
	It("Makes a http get request to each url in a file", func() {
		port := 8000
		list := `http://127.0.0.1:8000/A
http://127.0.0.1:8000/B
http://127.0.0.1:8000/C
		`

		urlsVisited := []string{}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlsVisited = append(urlsVisited, r.URL.Path)
		})

		server := httptest.NewUnstartedServer(handler)
		server.Listener, _ = net.Listen("tcp", ":"+strconv.Itoa(port))
		server.Start()
		defer server.Close()
		file, err := ioutil.TempFile(os.TempDir(), "prefix")
		if err != nil {
			panic(err)
		}
		file.WriteString(list)
		defer os.Remove(file.Name())

		log.Printf(fmt.Sprintf("Temp filename = %s", file.Name()))
		exePath, err := filepath.Abs("./code-named-something")
		if err != nil {
			panic(err)
		}
		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(fmt.Sprintf("OUTPUT : %v\n ERROR: %v" + string(output), err))

		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/A", nil)
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("%v", err)
			panic(err)
		}

		Expect(err).To(BeNil())
		Expect(containsString(urlsVisited, "/A")).To(Equal(true))
		Expect(containsString(urlsVisited, "/B")).To(Equal(true))
		Expect(containsString(urlsVisited, "/C")).To(Equal(true))
	})
})
