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
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/test"
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

func ExecutePlanBuilder(planBuilder *yaml.PlanBuilder) ([]byte, error) {
	return test.ExecutePlanBuilder("./corcel", planBuilder)
}

func ExecutePlanFromData(plan string) ([]byte, error) {
	return test.ExecutePlanFromData("./corcel", plan)
}

func ExecutePlanFromDataForApplication(plan string) (statistics.AggregatorSnapShot, error) {
	return test.ExecutePlanFromDataForApplication("./corcel", plan, config.Configuration{})
}

func ExecuteListForApplication(list []string, configuration config.Configuration) (statistics.AggregatorSnapShot, error) {
	return test.ExecuteListForApplication(list, configuration)
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
