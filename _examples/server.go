package main

import (
	"github.com/aellwein/coap"
	"github.com/aellwein/slf4go"
	_ "github.com/aellwein/slf4go-zap-adaptor"
)

func main() {
	slf4go.GetLoggerFactory().SetDefaultLogLevel(slf4go.LevelDebug)
	slf4go.GetLoggerFactory().SetLoggingParameters(slf4go.LoggingParameters{
		"development": true,
	})
	logger := slf4go.GetLogger("server")

	server, err := coap.NewInsecureCoapServerWithDefaultParameters(coap.NewResource("/rd"))
	if err != nil {
		logger.Panic(err)
	}

	logger.Debug(server)

	err = server.Listen()

	if err != nil {
		logger.Panic(err)
	}
}
