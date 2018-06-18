package message

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"encoding/binary"
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
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 0x1337)

		Convey("When decoded", func() {
			_, err := Decode(b, nil)

			Convey("'Invalid Message Version' is indicated by error", func() {
				So(err, ShouldEqual, InvalidMessageVersion)
			})
		})
	})
}
