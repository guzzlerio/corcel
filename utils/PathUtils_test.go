package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	var dir string

	var createTestFile = func(tmpPath string) {
		os.MkdirAll(tmpPath, os.ModePerm)
		content := []byte("temporary file's content")
		testFilePath := filepath.Join(tmpPath, "temp.txt")
		if err := ioutil.WriteFile(testFilePath, content, 0666); err != nil {
			log.Fatal(err)
		}
	}

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "corcel")
		if err != nil {
			log.Fatal(err)
		}
	})

	AfterEach(func() {
		os.RemoveAll(dir) // clean up
	})

	It("Finds file at same depth", func() {
		targetPath := filepath.Join(dir, "1", "2")

		startingPath := targetPath

		os.MkdirAll(startingPath, os.ModePerm)

		wd, _ := os.Getwd()
		defer func() {
			os.Chdir(wd)
		}()
		os.Chdir(startingPath)

		createTestFile(targetPath)

		var result = FindFileUp("temp.txt")

		Expect(result).ToNot(Equal(""))
	})

	It("Finds file at first parent", func() {
		targetPath := filepath.Join(dir, "1")

		startingPath := filepath.Join(dir, "1", "2")

		os.MkdirAll(startingPath, os.ModePerm)

		wd, _ := os.Getwd()
		defer func() {
			os.Chdir(wd)
		}()
		os.Chdir(startingPath)

		createTestFile(targetPath)

		var result = FindFileUp("temp.txt")

		Expect(result).ToNot(Equal(""))
	})

	It("Finds file at second parent", func() {
		targetPath := filepath.Join(dir)

		startingPath := filepath.Join(dir, "1", "2")

		os.MkdirAll(startingPath, os.ModePerm)

		wd, _ := os.Getwd()
		defer func() {
			os.Chdir(wd)
		}()
		os.Chdir(startingPath)

		createTestFile(targetPath)

		var result = FindFileUp("temp.txt")

		Expect(result).ToNot(Equal(""))
	})

})
