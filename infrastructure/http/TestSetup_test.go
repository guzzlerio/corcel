package http_test

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/guzzlerio/rizo"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/logger"
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
