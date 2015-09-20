package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http/httptest"
	"net"
	"net/http"
	"strings"
	"os/exec"
	"io/ioutil"
	"os"
	"log"
	"fmt"
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
		list := `
		http://127.0.0.1:$port/A
		http://127.0.0.1:$port/B
		http://127.0.0.1:$port/C
		`

		urlsVisited := []string{}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("BANG")
			urlsVisited = append(urlsVisited, r.URL.Path)
		})

		server := httptest.NewUnstartedServer(handler)
		server.Listener, _ = net.Listen("tcp", ":"+strconv.Itoa(port))
		server.Start()
		defer server.Close()
		file, _ := ioutil.TempFile(os.TempDir(), "prefix")
		file.WriteString(strings.Replace(list,"$port", strconv.Itoa(port), -1))
		defer os.Remove(file.Name())

		log.Printf(fmt.Sprintf("Temp filename = %s", file.Name()))
		cmd := exec.Command("./code-named-something", "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(err)
		fmt.Println(string(output))

		Expect(err).To(BeNil())
		Expect(containsString(urlsVisited,"/A")).To(Equal(true))
		Expect(containsString(urlsVisited,"/B")).To(Equal(true))
		Expect(containsString(urlsVisited,"/C")).To(Equal(true))
		})
	})
