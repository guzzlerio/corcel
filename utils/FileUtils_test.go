package utils_test

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	. "github.com/guzzlerio/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Something struct {
	A string `json:"a"`
	B int    `json:"b"`
}

var _ = Describe("FileUtils", func() {

	It("CreateFileFromLines", func() {
		var lines = []string{"A", "B", "C"}

		var file = CreateFileFromLines(lines)

		data, err := ioutil.ReadFile(file.Name())

		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("A\nB\nC\n"))

	})

	It("MarshalYamlToFile", func() {
		var subject = Something{
			A: "something",
			B: 1,
		}

		marshalFile, marshalErr := MarshalYamlToFile(os.TempDir(), subject)

		Expect(marshalErr).To(BeNil())

		fileData, fileErr := ioutil.ReadFile(marshalFile.Name())

		Expect(fileErr).To(BeNil())

		var newSubject Something
		yaml.Unmarshal(fileData, &newSubject)

		Expect(newSubject.A).To(Equal("something"))
		Expect(newSubject.B).To(Equal(1))
	})

})
