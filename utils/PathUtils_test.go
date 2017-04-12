package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils(t *testing.T) {
	Convey("Utils", t, func() {

		var dir string

		var createTestFile = func(tmpPath string) {
			os.MkdirAll(tmpPath, os.ModePerm)
			content := []byte("temporary file's content")
			testFilePath := filepath.Join(tmpPath, "temp.txt")
			if err := ioutil.WriteFile(testFilePath, content, 0666); err != nil {
				log.Fatal(err)
			}
		}

		func() {
			var err error
			dir, err = ioutil.TempDir("", "corcel")
			if err != nil {
				log.Fatal(err)
			}
		}()

		defer func() {
			os.RemoveAll(dir) // clean up
		}()

		Convey("Finds file at same depth", func() {
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

			So(result, ShouldNotEqual, "")
		})

		Convey("Finds file at first parent", func() {
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

			So(result, ShouldNotEqual, "")
		})

		Convey("Finds file at second parent", func() {
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

			So(result, ShouldNotEqual, "")
		})

	})
}
