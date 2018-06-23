package message

import (
	"fmt"
	"github.com/aellwein/coap/util"
	c "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTooShortMessage(t *testing.T) {
	c.Convey("Given a too short message", t, func() {
		b := make([]byte, 2)

		c.Convey("When decoded", func() {
			_, err := DecodeMessage(b, nil)

			c.Convey("'Packet is too short' is indicated by error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestInvalidMessageVersion(t *testing.T) {
	c.Convey("Given a message with an invalid version", t, func() {
		b := []byte{0xCA, 0xFE, 0xBA, 0xBE}

		c.Convey("When decoded", func() {
			_, err := DecodeMessage(b, nil)

			c.Convey("'Invalid Message Version' is indicated by error", func() {
				c.So(err, c.ShouldEqual, InvalidMessageVersion)
			})
		})
	})
}

func TestInvalidTokenLength(t *testing.T) {
	c.Convey("Given a message with invalid token length", t, func() {
		b := []byte{0x4A, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		c.Convey("When decoded", func() {
			_, err := DecodeMessage(b, nil)

			c.Convey("'Invalid Token Length' is indicated by error", func() {
				c.So(err, c.ShouldEqual, InvalidTokenLength)
			})
		})
	})
}

func TestMessageIsCorruptedOnShortToken(t *testing.T) {
	c.Convey("Given a message with valid token length but short token content", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3}

		c.Convey("When decoded", func() {
			_, err := DecodeMessage(b, nil)

			c.Convey("'Packet Is Too Short' is indicated by error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestMessageCodeIsParsedCorrectly(t *testing.T) {
	c.Convey("Given a message with code class 0.2", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		c.Convey("When decoded", func() {
			msg, err := DecodeMessage(b, nil)

			c.Convey("Code should be 0.2", func() {
				c.So(err, c.ShouldBeNil)
				c.So(msg.Code.CodeClass, c.ShouldEqual, 0)
				c.So(msg.Code.CodeDetail, c.ShouldEqual, 2)
			})
		})
	})
}

func TestOptionsAreParsedCorrectly(t *testing.T) {
	c.Convey("Given a valid message with options", t, func() {
		/*
			00000000  44 02 56 b3 00 00 4a 82  39 6c 6f 63 61 6c 68 6f  |D.V...J.9localho|
			00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
			00000020  03 62 3d 55 06 6c 74 3d  33 30 30                 |.b=U.lt=300|
		*/
		b := []byte{
			0x44, 0x02, 0x56, 0xB3, 0x00, 0x00, 0x4A, 0x82, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C,
			0x68, 0x6F, 0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D,
			0x61, 0x6C, 0x65, 0x78, 0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30,
			0x30,
		}
		c.Convey("When decoded", func() {
			m, err := DecodeMessage(b, nil)

			c.Convey("Decode should not have errors", func() {
				c.So(err, c.ShouldBeNil)
				fmt.Printf("%v", m)
			})
			c.Convey("Message should have 4 options", func() {
				c.So(len(*m.Options), c.ShouldEqual, 4)
			})
			c.Convey("Message should have Uri-Host to be set to 'localhost'", func() {
				v, ok := (*m.Options)[UriHost]
				c.So(ok, c.ShouldBeTrue)
				c.So(string(v[0]), c.ShouldEqual, "localhost")
			})
			c.Convey("Message should have Uri-Path to be set to 'rd'", func() {
				v, ok := (*m.Options)[UriPath]
				c.So(ok, c.ShouldBeTrue)
				c.So(string(v[0]), c.ShouldEqual, "rd")
			})
			c.Convey("Message should have Uri-Port to be set to 5683", func() {
				v, ok := (*m.Options)[UriPort]
				c.So(ok, c.ShouldBeTrue)
				c.So(util.ToBigEndianNumber(v[0]).(uint16), c.ShouldEqual, 5683)
			})

		})
	})
}
