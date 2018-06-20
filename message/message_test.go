package message

import (
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
func TestMessageToString(t *testing.T) {
	c.Convey("Given a valid message", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		c.Convey("When decoded", func() {
			msg, err := DecodeMessage(b, nil)

			c.Convey("the stringified message must be equal to expected output", func() {
				c.So(err, c.ShouldBeNil)
				c.So(msg.String(), c.ShouldEqual, "Message{type=CON, code=0.02, id=8818, tkn=0x0471BD4AF3A34709, options=&map[], payload=[0x02, 0x22, 0x72, 0x04, 0x71, 0xBD, 0x4A, 0xF3, 0xA3, 0x47, 0x09], from=<nil>}")
			})
		})
	})
}

func TestParseOptions(t *testing.T) {
	c.Convey("Given a valid message with options", t, func() {
		b := []byte{
			0x44, 0x02, 0x16, 0xb3, 0x00, 0x00, 0xe3, 0x99, 0x33, 0x3a, 0x3a, 0x31, 0x42, 0x16, 0x33, 0x42, 0x72,
			0x64, 0x47, 0x65, 0x70, 0x3d, 0x61, 0x6c, 0x65, 0x78, 0x03, 0x62, 0x3d, 0x55, 0x06, 0x6c, 0x74,
			0x3d, 0x33, 0x30, 0x30,
		}
		c.Convey("When decoded", func() {
			_, err := DecodeMessage(b, nil)

			c.Convey("Code should be 0.2", func() {
				c.So(err, c.ShouldBeNil)
			})
		})
	})
}
