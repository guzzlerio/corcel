package utils_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	. "github.com/guzzlerio/corcel/utils"

	. "github.com/smartystreets/goconvey/convey"
)

type Something struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestFileUtils(t *testing.T) {
	Convey("FileUtils", t, func() {

		Convey("CreateFileFromLines", func() {
			var lines = []string{"A", "B", "C"}

			var file = CreateFileFromLines(lines)

			data, err := ioutil.ReadFile(file.Name())

			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, "A\nB\nC\n")

		})

		Convey("MarshalYamlToFile", func() {
			var subject = Something{
				A: "something",
				B: 1,
			}

			marshalFile, marshalErr := MarshalYamlToFile(os.TempDir(), subject)

			So(marshalErr, ShouldBeNil)

			fileData, fileErr := ioutil.ReadFile(marshalFile.Name())

			So(fileErr, ShouldBeNil)

			var newSubject Something
			yaml.Unmarshal(fileData, &newSubject)

			So(newSubject.A, ShouldEqual, "something")
			So(newSubject.B, ShouldEqual, 1)
		})

	})
}
