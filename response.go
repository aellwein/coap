package coap

import (
	"fmt"
	"github.com/aellwein/coap/message"
)

// Response represents a response message, which always follows a CON request.
type Response interface {
	GetResp() *message.Message
}

// internal responseAdapter type
type responseAdapter struct {
	resp *message.Message
}

func (r *responseAdapter) String() string {
	return fmt.Sprintf("Response{ code: %v (%v) }", r.resp.Code, r.resp.Code.ToResponseCode())
}

func (r *responseAdapter) GetResp() *message.Message {
	return r.resp
}

func newResponseAdapter(code *message.ResponseCode, request *Request) Response {
	c := code.ToCodeType()
	msg := message.NewAcknowledgementMessage(c)
	msg.Token = (*request).GetReq().Token // important to use the token from request
	return &responseAdapter{resp: msg}
}

/* Success responses */

func NewOkResponse(request *Request) Response {
	return newResponseAdapter(message.Ok, request)
}

func NewCreatedResponse(request *Request) Response {
	return newResponseAdapter(message.Created, request)
}
func NewDeletedResponse(request *Request) Response {
	return newResponseAdapter(message.Deleted, request)
}
func NewValidResponse(request *Request) Response {
	return newResponseAdapter(message.Valid, request)
}
func NewChangedResponse(request *Request) Response {
	return newResponseAdapter(message.Changed, request)
}
func NewContentResponse(request *Request) Response {
	return newResponseAdapter(message.Content, request)
}

/* Client error responses */

func NewBadRequestResponse(request *Request) Response {
	return newResponseAdapter(message.BadRequest, request)
}

func NewUnauthorizedResponse(request *Request) Response {
	return newResponseAdapter(message.Unauthorized, request)
}
func NewBadOptionResponse(request *Request) Response {
	return newResponseAdapter(message.BadOption, request)
}

func NewForbiddenResponse(request *Request) Response {
	return newResponseAdapter(message.Forbidden, request)
}
func NewNotFoundResponse(request *Request) Response {
	return newResponseAdapter(message.NotFound, request)
}
func NewMethodNotAllowedResponse(request *Request) Response {
	return newResponseAdapter(message.MethodNotAllowed, request)
}
func NewNotAcceptableResponse(request *Request) Response {
	return newResponseAdapter(message.NotAcceptable, request)
}
func NewPreconditionFailedResponse(request *Request) Response {
	return newResponseAdapter(message.PreconditionFailed, request)
}
func NewRequestEntityTooLargeResponse(request *Request) Response {
	return newResponseAdapter(message.RequestEntityTooLarge, request)
}
func NewUnsupportedContentFormatResponse(request *Request) Response {
	return newResponseAdapter(message.UnsupportedContentFormat, request)
}

/* server error codes */

func NewInternalServerErrorResponse(request *Request) Response {
	return newResponseAdapter(message.InternalServerError, request)
}

func NewNotImplementedResponse(request *Request) Response {
	return newResponseAdapter(message.NotImplemented, request)
}
func NewBadGatewayResponse(request *Request) Response {
	return newResponseAdapter(message.BadGateway, request)
}
func NewServiceUnavailableResponse(request *Request) Response {
	return newResponseAdapter(message.ServiceUnavailable, request)
}
func NewGatewayTimeoutResponse(request *Request) Response {
	return newResponseAdapter(message.GatewayTimeout, request)
}
func NewProxyingNotSupportedResponse(request *Request) Response {
	return newResponseAdapter(message.ProxyingNotSupported, request)
}
