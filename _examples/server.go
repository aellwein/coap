package main

import (
	"github.com/aellwein/coap"
	"github.com/aellwein/slf4go"
	_ "github.com/aellwein/slf4go-zap-adaptor"
)

var (
	server *coap.Server
	logger slf4go.Logger
)

func init() {
	slf4go.GetLoggerFactory().SetDefaultLogLevel(slf4go.LevelDebug)
	logger = slf4go.GetLogger("server")
}

func onRegister(request *coap.Message) (*coap.Message, error) {

	logger.Debugf("onRegister(%v)", request)

	server.AddResource(&coap.Resource{
		Path:     "/rd/cafe/babe",
		OnPOST:   OnUpdate,
		OnDELETE: OnDelete,
	})

	return coap.NewAcknowledgementMessageBuilder().
		Code(coap.Created).
		MessageId(request.MessageID).
		Token(request.Token).
		Option(coap.LocationPath, coap.NewLocationPathOption("/rd/cafe/babe")...).
		Build(), nil

}

func OnUpdate(request *coap.Message) (*coap.Message, error) {

	logger.Debugf("onUpdate(%v)", request)

	return coap.NewAcknowledgementMessageBuilder().
		Code(coap.Changed).
		MessageId(request.MessageID).
		Token(request.Token).
		Option(coap.LocationPath, coap.NewLocationPathOption("/rd/cafe/babe")...).
		Build(), nil
}

func OnDelete(request *coap.Message) (*coap.Message, error) {

	logger.Debugf("onDelete(%v)", request)

	server.RemoveResourceByPath("/rd/cafe/babe")

	return coap.NewAcknowledgementMessageBuilder().
		Code(coap.Deleted).
		MessageId(request.MessageID).
		Token(request.Token).
		Option(coap.LocationPath, coap.NewLocationPathOption("/rd/cafe/babe")...).
		Build(), nil
}

func main() {
	var err error
	server, err = coap.NewInsecureCoapServerWithDefaultParameters(
		&coap.Resource{
			Path:   "/rd",
			OnPOST: onRegister,
		},
	)
	if err != nil {
		logger.Panic(err)
	}

	logger.Debug(server)

	err = server.Listen()

	if err != nil {
		logger.Panic(err)
	}
}
