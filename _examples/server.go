package main

import (
	"fmt"
	"github.com/aellwein/coap"
	"github.com/aellwein/slf4go"
	_ "github.com/aellwein/slf4go-zap-adaptor"
)

func onPOST(request coap.Request) error {
	fmt.Println(request)
	return nil
}

func main() {
	slf4go.GetLoggerFactory().SetLoggingParameters(slf4go.LoggingParameters{
		"development": true,
	})

	logger := slf4go.GetLogger("server")
	logger.SetLevel(slf4go.LevelTrace)

	server, err := coap.NewInsecureCoapServerWithDefaultParameters(coap.RequestHandler{Path: "/rd", HandlePOST: onPOST})
	if err != nil {
		logger.Panic(err)
	}

	logger.Debug(server)

	err = server.Listen()

	if err != nil {
		logger.Panic(err)
	}
}
