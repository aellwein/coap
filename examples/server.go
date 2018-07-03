package main

import (
	"fmt"
	"github.com/aellwein/coap"
	"github.com/aellwein/coap/logging"
)

func onPOST(request coap.Request) error {
	fmt.Println(request)
	return nil
}

func main() {
	logger := logging.LoggerFactory.GetLogger("server")
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
