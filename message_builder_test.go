package coap

import (
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
				c.So(msg.Payload, c.ShouldResemble, PayloadType("lalala"))
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
				WithPayload(ContentTypeApplicationJson, []byte("lalala")).
				Build()

			c.Convey("Then all of the fields are set correctly", func() {
				c.So(*msg.Code, c.ShouldResemble, *GET)
				c.So(msg.MessageID, c.ShouldNotEqual, 0)
				c.So(msg.Token, c.ShouldNotResemble, []byte{0, 0, 0, 0, 0, 0, 0, 0})
				c.So(len(*msg.Options), c.ShouldEqual, 2)
				c.So((*msg.Options)[UriPath], c.ShouldResemble, []OptionValueType{[]byte("rd")})
				c.So(msg.Payload, c.ShouldResemble, PayloadType("lalala"))
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
				c.So(msg.Payload, c.ShouldBeNil)
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
				c.So(msg.Payload, c.ShouldBeNil)
			})
		})
	})

}
