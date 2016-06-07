package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/global"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/test"
	"github.com/Sirupsen/logrus"
	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	//TestServer ...
	TestServer *rizo.RequestRecordingServer
)

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Corcel Suite")
}

func ExecutePlanBuilder(planBuilder *test.YamlPlanBuilder) error {
	file, err := planBuilder.Build()
	if err != nil {
		return err
	}
	exePath, err := filepath.Abs("./corcel")
	if err != nil {
		return err
	}
	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()
	args := []string{"--plan"}
	cmd := exec.Command(exePath, append(args, file.Name())...)
	output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	logger.Log.Println(fmt.Sprintf("%s", output))
	return err
}

var _ = BeforeSuite(func() {
	logger.Initialise()
	logger.ConfigureLogging(&config.Configuration{})
	logrus.SetOutput(ioutil.Discard)
	logger.Log.Out = ioutil.Discard
	TestServer = rizo.CreateRequestRecordingServer(global.TestPort)
	TestServer.Start()
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})
