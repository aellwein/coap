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
	logging.InitLogger(true)
	server, err := coap.NewInsecureCoapServer(coap.RequestHandler{Path: "/rd", HandlePOST: onPOST})

	if err != nil {
		logging.Sugar.Panic(err)
	}

	logging.Sugar.Debug(server)

	err = server.Listen()

	if err != nil {
		logging.Sugar.Panic(err)
	}
}
