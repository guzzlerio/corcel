package request

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/logger"

	"github.com/guzzlerio/rizo"
)

var (
	//TestServer ...
	TestServer *rizo.RequestRecordingServer
)

func BeforeTest() {
	logger.Initialise()
	logger.ConfigureLogging(&config.Configuration{})
	logrus.SetOutput(ioutil.Discard)
	logger.Log.Out = ioutil.Discard
	TestServer = rizo.CreateRequestRecordingServer(global.TestPort)
	TestServer.Start()
}

func AfterTest() {
	TestServer.Stop()
}
