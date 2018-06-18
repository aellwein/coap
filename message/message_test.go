package message

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTooShortMessage(t *testing.T) {
	Convey("Given a too short message", t, func() {
		b := make([]byte, 2)

		Convey("When decoded", func() {
			_, err := Decode(b, nil)

			Convey("'Packet is too short' is indicated by error", func() {
				So(err, ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestInvalidMessageVersion(t *testing.T) {
	Convey("Given a message with an invalid version", t, func() {
		b := []byte{0xCA, 0xFE, 0xBA, 0xBE}

		Convey("When decoded", func() {
			_, err := Decode(b, nil)

			Convey("'Invalid Message Version' is indicated by error", func() {
				So(err, ShouldEqual, InvalidMessageVersion)
			})
		})
	})
}

func TestInvalidTokenLength(t *testing.T) {
	Convey("Given a message with invalid token length", t, func() {
		b := []byte{0x4A, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		Convey("When decoded", func() {
			_, err := Decode(b, nil)

			Convey("'Invalid Token Length' is indicated by error", func() {
				So(err, ShouldEqual, InvalidTokenLength)
			})
		})
	})
}

func TestMessageIsCorruptedOnShortToken(t *testing.T) {
	Convey("Given a message with valid token length but short token content", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3}

		Convey("When decoded", func() {
			_, err := Decode(b, nil)

			Convey("'Packet Is Too Short' is indicated by error", func() {
				So(err, ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestMessageCodeIsParsedCorrectly(t *testing.T) {
	Convey("Given a message with code class 0.2", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		Convey("When decoded", func() {
			msg, err := Decode(b, nil)

			Convey("Code should be 0.2", func() {
				So(err, ShouldBeNil)
				So(msg.Code.CodeClass, ShouldEqual, 0)
				So(msg.Code.CodeDetail, ShouldEqual, 2)
			})
		})
	})
}
func TestMessageToString(t *testing.T) {
	Convey("Given a valid message", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		Convey("When decoded", func() {
			msg, err := Decode(b, nil)

			Convey("the stringified message must be equal to expected output", func() {
				So(err, ShouldBeNil)
				So(msg.String(), ShouldEqual, "Message{type=CON, code=0.02, id=8818, tkn=0x471BD4AF3A34709 (8), from=<nil>}")
			})
		})
	})
}

