package main

import (
	"io/ioutil"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/rizo"
)

var (
	//TestServer ...
	TestServer *rizo.RequestRecordingServer
)

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Corcel Suite")
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
