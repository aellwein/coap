package coap

import (
	"fmt"
	"github.com/aellwein/coap/message"
)

type Request interface {
	Message() message.Message
}

type RequestHandler struct {
	Path         string
	HandlePOST   func(r Request) error
	HandlePUT    func(r Request) error
	HandleGET    func(r Request) error
	HandleDELETE func(r Request) error
}

func (r RequestHandler) String() string {
	return fmt.Sprintf("RequestHandler{ path=%s, HandlePOST=%p, HandlePUT=%p, HandleGET=%p, HandleDELETE=%p }",
		r.Path,
		r.HandlePOST,
		r.HandlePUT,
		r.HandleGET,
		r.HandleDELETE)
}
