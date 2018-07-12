package coap

import (
	"fmt"
	"github.com/aellwein/coap/message"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

type mockRequest struct {
	msg *message.Message
}

func (r *mockRequest) GetReq() *message.Message {
	return r.msg
}

func newMockRequest() Request {
	b := []byte{
		0x44, 0x02, 0x1B, 0x2B, 0x00, 0x00, 0x3F, 0x3D, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
		0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
		0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0B, 0x30, 0x31, 0x32, 0x33,
		0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
	}
	m, _ := message.Decode(b, nil)
	return &mockRequest{msg: m}
}

type responses struct {
	code               message.ResponseCode
	responseCreateFunc func(r *Request) Response
}

var allResponses = []responses{
	{
		code:               message.Created,
		responseCreateFunc: NewCreatedResponse,
	},
	{
		code:               message.Deleted,
		responseCreateFunc: NewDeletedResponse,
	},
	{
		code:               message.Valid,
		responseCreateFunc: NewValidResponse,
	},
	{
		code:               message.Changed,
		responseCreateFunc: NewChangedResponse,
	},
	{
		code:               message.Content,
		responseCreateFunc: NewContentResponse,
	},
	{
		code:               message.BadRequest,
		responseCreateFunc: NewBadRequestResponse,
	},
	{
		code:               message.Unauthorized,
		responseCreateFunc: NewUnauthorizedResponse,
	},
	{
		code:               message.BadOption,
		responseCreateFunc: NewBadOptionResponse,
	},
	{
		code:               message.Forbidden,
		responseCreateFunc: NewForbiddenResponse,
	},
	{
		code:               message.NotFound,
		responseCreateFunc: NewNotFoundResponse,
	},
	{
		code:               message.MethodNotAllowed,
		responseCreateFunc: NewMethodNotAllowedResponse,
	},
	{
		code:               message.NotAcceptable,
		responseCreateFunc: NewNotAcceptableResponse,
	},
	{
		code:               message.PreconditionFailed,
		responseCreateFunc: NewPreconditionFailedResponse,
	},
	{
		code:               message.RequestEntityTooLarge,
		responseCreateFunc: NewRequestEntityTooLargeResponse,
	},
	{
		code:               message.UnsupportedContentFormat,
		responseCreateFunc: NewUnsupportedContentFormatResponse,
	},
	{
		code:               message.InternalServerError,
		responseCreateFunc: NewInternalServerErrorResponse,
	},
	{
		code:               message.NotImplemented,
		responseCreateFunc: NewNotImplementedResponse,
	},
	{
		code:               message.BadGateway,
		responseCreateFunc: NewBadGatewayResponse,
	},
	{
		code:               message.ServiceUnavailable,
		responseCreateFunc: NewServiceUnavailableResponse,
	},
	{
		code:               message.GatewayTimeout,
		responseCreateFunc: NewGatewayTimeoutResponse,
	},
	{
		code:               message.ProxyingNotSupported,
		responseCreateFunc: NewProxyingNotSupportedResponse,
	},
}

func TestNewResponse_Auto(t *testing.T) {
	convey.Convey("Given a mock request", t, func() {
		mr := newMockRequest()

		for _, r := range allResponses {
			convey.Convey(fmt.Sprintf("When response '%v' is created", r.code), func() {
				resp := r.responseCreateFunc(&mr)

				convey.Convey("Then response code matches the set up code", func() {
					convey.So(*resp.GetResp().Code, convey.ShouldResemble, r.code.ToCodeType())

					convey.Convey("And token must be equal the request token", func() {

						convey.So(resp.GetResp().Token, convey.ShouldResemble, mr.GetReq().Token)
					})
				})
			})
		}
	})
}
