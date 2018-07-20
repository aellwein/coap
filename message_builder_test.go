package coap

import (
	"encoding/binary"
	"fmt"
	c "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestNewConfirmableMessageBuilder(t *testing.T) {
	c.Convey("Given a new message builder", t, func() {
		c.Convey("When a CON message is created", func() {

			msg := NewConfirmableMessageBuilder().
				Code(POST).
				WithRandomMessageId().
				WithRandomToken().
				Option(UriPath, []byte("rd")).
				Option(UriQuery, []byte{}).
				WithPayload(ContentTypeApplicationJson, []byte("lalala")).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *POST)
				c.So(msg.MessageID, c.ShouldNotEqual, 0)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0, 0, 0, 0, 0, 0, 0, 0})
				c.So(len(*msg.Options), c.ShouldEqual, 3)
				c.So((*msg.Options)[UriPath], c.ShouldResemble, []OptionValueType{[]byte("rd")})
				c.So(*msg.Payload.Type, c.ShouldEqual, ContentTypeApplicationJson)
				c.So(msg.Payload.Content, c.ShouldResemble, []byte("lalala"))
			})
		})
	})
}

func TestNewConfirmableMessageBuilderWithCustomPayloadType(t *testing.T) {
	c.Convey("Given a new message builder", t, func() {
		c.Convey("When a CON message is created", func() {

			msg := NewConfirmableMessageBuilder().
				Code(POST).
				WithRandomMessageId().
				WithRandomToken().
				Option(UriPath, []byte("rd")).
				Option(UriQuery, []byte{}).
				WithPayload(65000, []byte("lalala")).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *POST)
				c.So(msg.MessageID, c.ShouldNotEqual, 0)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0, 0, 0, 0, 0, 0, 0, 0})
				c.So(len(*msg.Options), c.ShouldEqual, 3)
				c.So((*msg.Options)[UriPath], c.ShouldResemble, []OptionValueType{[]byte("rd")})
				c.So(binary.BigEndian.Uint16((*msg.Options)[ContentFormat][0]), c.ShouldEqual, 65000)
				c.So(*msg.Payload.Type, c.ShouldEqual, 65000)
				c.So(msg.Payload.Content, c.ShouldResemble, []byte("lalala"))
			})
		})
	})
}

func TestNewNonConfirmableMessageBuilder(t *testing.T) {
	c.Convey("Given a new message builder", t, func() {
		c.Convey("When a NON message is created", func() {
			addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1")
			msg := NewNonConfirmableMessageBuilder().
				From(addr).
				Code(GET).
				WithRandomMessageId().
				WithRandomToken().
				Option(UriPath, []byte("rd")).
				WithPayload(ContentTypeTextPlain, []byte("lalala")).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *GET)
				c.So(msg.MessageID, c.ShouldNotEqual, 0)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0, 0, 0, 0, 0, 0, 0, 0})
				c.So(len(*msg.Options), c.ShouldEqual, 2)
				c.So((*msg.Options)[UriPath], c.ShouldResemble, []OptionValueType{[]byte("rd")})
				c.So(*msg.Payload.Type, c.ShouldEqual, ContentTypeTextPlain)
				c.So(msg.Payload.Content, c.ShouldResemble, []byte("lalala"))
				c.So(msg.Source, c.ShouldEqual, addr)
			})
		})
	})
}

func TestNewAcknowledgementMessageBuilder(t *testing.T) {
	c.Convey("Given a new message builder", t, func() {
		c.Convey("When a ACK message is created", func() {
			tkn := TokenType{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF}
			msg := NewAcknowledgementMessageBuilder().
				Code(Created).
				MessageId(0x1337).
				Token(&tkn).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *Created)
				c.So(msg.MessageID, c.ShouldEqual, 0x1337)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF})
				c.So(len(*msg.Options), c.ShouldEqual, 0)
				c.So(msg.Payload.Type, c.ShouldBeNil)
				c.So(msg.Payload.Content, c.ShouldBeNil)
			})
		})
	})
}

func TestNewResetMessageBuilder(t *testing.T) {
	c.Convey("Given a new message builder", t, func() {
		c.Convey("When a RST message is created", func() {
			tkn := TokenType{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF}
			msg := NewResetMessageBuilder().
				Code(Created).
				MessageId(0x1337).
				Token(&tkn).
				Option(LocationPath, NewLocationPathOption("/cafe")...).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *Created)
				c.So(msg.MessageID, c.ShouldEqual, 0x1337)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF})
				c.So(len(*msg.Options), c.ShouldEqual, 1)
				c.So(msg.Payload.Type, c.ShouldBeNil)
				c.So(msg.Payload.Content, c.ShouldBeNil)
			})
		})
	})
}

func TestNewMessageFromBytesAndPeer(t *testing.T) {
	c.Convey("Given a message as byte array", t, func() {
		b := []byte{
			0x44, 0x02, 0x1B, 0x2B, 0x00, 0x00, 0x3F, 0x3D, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0B, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xFF, 0x61, 0x6C, 0x65, 0x78, 0x31, 0x32, 0x33,
		}
		peer := net.UDPAddr{Port: 5683, IP: []byte{127, 0, 0, 1}}

		c.Convey("When a message is created and the source IP is set", func() {
			msg, err := NewMessageFromBytesAndPeer(b, &peer)

			c.Convey("Then the result message is valid and contains source IP", func() {
				c.So(err, c.ShouldBeNil)
				c.So(*msg.Source, c.ShouldResemble, peer)
			})
		})
	})
}

type messagesTest struct {
	code       *CodeType
	createFunc func(message *Message) *Message
}

func TestNewResponseMessage_Auto(t *testing.T) {
	var msgs = []messagesTest{
		{
			code:       BadRequest,
			createFunc: NewBadRequestResponseMessage,
		},
		{
			code:       Unauthorized,
			createFunc: NewUnauthorizedResponseMessage,
		},
		{
			code:       BadOption,
			createFunc: NewBadOptionResponseMessage,
		},
		{
			code:       Forbidden,
			createFunc: NewForbiddenResponseMessage,
		},
		{
			code:       NotFound,
			createFunc: NewNotFoundResponseMessage,
		},
		{
			code:       MethodNotAllowed,
			createFunc: NewMethodNotAllowedResponseMessage,
		},
		{
			code:       NotAcceptable,
			createFunc: NewNotAcceptableResponseMessage,
		},
		{
			code:       PreconditionFailed,
			createFunc: NewPreconditionFailedResponseMessage,
		},
		{
			code:       RequestEntityTooLarge,
			createFunc: NewRequestEntityTooLargeResponseMessage,
		},
		{
			code:       UnsupportedContentFormat,
			createFunc: NewUnsupportedContentFormatResponseMessage,
		},
		{
			code:       InternalServerError,
			createFunc: NewInternalServerErrorResponseMessage,
		},
		{
			code:       NotImplemented,
			createFunc: NewNotImplementedResponseMessage,
		},
		{
			code:       BadGateway,
			createFunc: NewBadGatewayResponseMessage,
		},
		{
			code:       ServiceUnavailable,
			createFunc: NewServiceUnavailableResponseMessage,
		},
		{
			code:       GatewayTimeout,
			createFunc: NewGatewayTimeoutResponseMessage,
		},
		{
			code:       ProxyingNotSupported,
			createFunc: NewProxyingNotSupportedResponseMessage,
		},
	}

	c.Convey("Given a new request message", t, func() {
		req, _ := NewMessageFromBytes([]byte{
			0x44, 0x02, 0x1B, 0x2B, 0x00, 0x00, 0x3F, 0x3D, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0B, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xFF, 0x61, 0x6C, 0x65, 0x78, 0x31, 0x32, 0x33,
		})

		for _, tm := range msgs {
			c.Convey(fmt.Sprintf("When a new response message with response type %v is created", tm.code), func() {
				m := tm.createFunc(req)

				expected := tm.code

				c.Convey(fmt.Sprintf("Then the response message contains the expected code %v", expected), func() {
					c.So(*m.Code, c.ShouldResemble, *expected)
					c.Convey("And the message ID is equal to the message ID of the request message", func() {
						c.So(m.MessageID, c.ShouldEqual, req.MessageID)
						c.Convey("And token of the response message is equal to token of the request message", func() {
							c.So(*m.Token, c.ShouldResemble, *req.Token)
						})
					})
				})
			})
		}
	})
}
