package http_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/guzzlerio/rizo"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/test"
)

var (
	//TestServer ...
	TestServer *rizo.RequestRecordingServer
)

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Suite")
}

var _ = BeforeSuite(func() {
	logger.Initialise()
	logger.ConfigureLogging(&config.Configuration{})
	logrus.SetOutput(ioutil.Discard)
	logger.Log.Out = ioutil.Discard
	TestServer = rizo.CreateRequestRecordingServer(global.TestPort)
	TestServer.Start()

	os.Remove("./output.yml")
	os.Remove("./history.yml")
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})

func ExecutePlanFromData(plan string) ([]byte, error) {
	return test.ExecutePlanFromData("../.././corcel", plan)
}

func ExecutePlanFromDataForApplication(plan string) (statistics.AggregatorSnapShot, error) {
	return test.ExecutePlanFromDataForApplication("../.././corcel", plan, config.Configuration{})
}

func ExecutePlanBuilder(planBuilder *yaml.PlanBuilder) ([]byte, error) {
	return test.ExecutePlanBuilder("../.././corcel", planBuilder)
}

func ExecutePlanBuilderForApplication(planBuilder *yaml.PlanBuilder) (statistics.AggregatorSnapShot, error) {
	return test.ExecutePlanBuilderForApplication("../.././corcel", planBuilder, config.Configuration{})
}
