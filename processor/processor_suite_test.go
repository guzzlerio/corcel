package processor

import (
	"io/ioutil"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/logger"
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

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Processor Suite")
}
