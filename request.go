package coap

import "github.com/aellwein/coap/message"

// Request describes an incoming request.
type Request interface {
	GetReq() *message.Message
}
