package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/satori/go.uuid"
)

//CreateFileFromLines ...
func CreateFileFromLines(lines []string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	errormanager.Check(err)
	for _, line := range lines {
		_, err := file.WriteString(fmt.Sprintf("%s\n", line))
		errormanager.Check(err)
	}
	errormanager.Check(file.Sync())
	return file
}

//MarshalYamlToFile
func MarshalYamlToFile(outputPath string, object interface{}) (*os.File, error) {
	contents, err := yaml.Marshal(&object)
	if err != nil {
		return nil, err
	}

	file, err := ioutil.TempFile(os.TempDir(), "yamlExecutionPlanForCorcel")
	if err != nil {
		return nil, err
	}
	defer func() {
		CheckErr(file.Close())
	}()
	//FIXME Write returns an error which is ignored...
	file.Write(contents)

	//FIXME ignored error output from MkdirAll
	os.MkdirAll(outputPath, 0777)

	err = ioutil.WriteFile(path.Join(outputPath, uuid.NewV4().String()), contents, 0644)
	if err != nil {
		panic(err)
	}

	err = file.Sync()

	if err != nil {
		return nil, err
	}
	return file, nil
}
