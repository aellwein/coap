package message

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTooShortMessageCausesError(t *testing.T) {
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

func TestInvalidMessageVersionCausesError(t *testing.T) {
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

func TestMessageCodeIsParsedCorrectly(t *testing.T) {
	Convey("Given a message with code class ", t, func() {
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
