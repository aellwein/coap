package message

import (
	"fmt"
	"github.com/aellwein/coap/util"
	c "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestTooShortMessage(t *testing.T) {
	c.Convey("Given a too short message", t, func() {
		b := make([]byte, 2)

		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then 'Packet is too short' is indicated by error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestInvalidMessageVersion(t *testing.T) {
	c.Convey("Given a message with an invalid version", t, func() {
		b := []byte{0xCA, 0xFE, 0xBA, 0xBE}

		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then 'Invalid Message Version' is indicated by error", func() {
				c.So(err, c.ShouldEqual, InvalidMessageVersion)
			})
		})
	})
}

func TestInvalidTokenLength(t *testing.T) {
	c.Convey("Given a message with invalid token length", t, func() {
		b := []byte{0x4A, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then 'Invalid Token Length' is indicated by error", func() {
				c.So(err, c.ShouldEqual, InvalidTokenLength)
			})
		})
	})
}

func TestMessageIsCorruptedOnShortToken(t *testing.T) {
	c.Convey("Given a message with valid token length but short token content", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3}

		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then 'Packet Is Too Short' is indicated by error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestMessageCodeIsParsedCorrectly(t *testing.T) {
	c.Convey("Given a message with code class 0.2", t, func() {
		b := []byte{0x48, 0x02, 0x22, 0x72, 0x04, 0x71, 0xbd, 0x4a, 0xf3, 0xa3, 0x47, 0x09}

		c.Convey("When decoded", func() {
			msg, err := Decode(b, nil)

			c.Convey("Then code should be 0.2", func() {
				c.So(err, c.ShouldBeNil)
				c.So(msg.Code.CodeClass, c.ShouldEqual, 0)
				c.So(msg.Code.CodeDetail, c.ShouldEqual, 2)
			})
		})
	})
}

func TestOptionsAreParsedCorrectly(t *testing.T) {
	/*
		00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
		00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
		00000020  03 62 3d 55 06 6c 74 3d  33 30 30 0d ed 30 31 32  |.b=U.lt=300..012|
		00000030  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000040  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		00000050  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000060  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000070  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000080  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000090  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000a0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		000000b0  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		000000c0  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		000000d0  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		000000e0  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000f0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000100  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000110  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000120  33 34 35 36 37 38 39                              |3456789|
	*/
	c.Convey("Given a valid message with options", t, func() {
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0D, 0xED, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		}
		c.Convey("When decoded", func() {
			m, err := Decode(b, nil)

			c.Convey("Then decode should not have errors", func() {
				c.So(err, c.ShouldBeNil)
				fmt.Printf("%v", m)
			})
			c.Convey("And message should have 4 options", func() {
				c.So(len(*m.Options), c.ShouldEqual, 4)
			})
			c.Convey("And message should have Uri-Host to be set to 'localhost'", func() {
				v, ok := (*m.Options)[UriHost]
				c.So(ok, c.ShouldBeTrue)
				c.So(string(v[0]), c.ShouldEqual, "localhost")
			})
			c.Convey("And message should have Uri-Path to be set to 'rd'", func() {
				v, ok := (*m.Options)[UriPath]
				c.So(ok, c.ShouldBeTrue)
				c.So(string(v[0]), c.ShouldEqual, "rd")
			})
			c.Convey("And message should have Uri-Port to be set to 5683", func() {
				v, ok := (*m.Options)[UriPort]
				c.So(ok, c.ShouldBeTrue)
				c.So(util.ToBigEndianNumber(v[0]).(uint16), c.ShouldEqual, 5683)
			})
			c.Convey("And 3rd Uri-Query parameter should be equal a set string", func() {
				expected := strings.Join(
					[]string{
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
					}, "")

				v, ok := (*m.Options)[UriQuery]
				c.So(ok, c.ShouldBeTrue)
				c.ShouldEqual(string(v[3]), expected)
			})
		})
	})
}

func TestInvalidOptionLengthCausesError(t *testing.T) {
	/*
		00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
		00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
		00000020  03 62 3d 55 06 6c 74 3d  33 30 30 0d ed 30 31 32  |.b=U.lt=300..012|
		00000030  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000040  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		00000050  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000060  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000070  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000080  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000090  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000a0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		000000b0  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		000000c0  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		000000d0  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		000000e0  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000f0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000100  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000110  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000120  33 34 35 36 37 38 39                              |3456789|
	*/
	c.Convey("Given a message with one invalid option length", t, func() {
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0D, 0xF0, 0x30, 0x31, 0x32, // 0xF0: manipulated
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
			0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
			0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		}
		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then decode should fail with 'PacketIsTooShort' error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestMissingOptionExtendedValueCausesError(t *testing.T) {
	/*
		00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
		00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
		00000020  03 62 3d 55 06 6c 74 3d  33 30 30 0d ed 30 31 32  |.b=U.lt=300..012|
		00000030  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000040  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		00000050  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000060  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000070  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000080  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		00000090  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000a0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		000000b0  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		000000c0  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		000000d0  33 34 35 36 37 38 39 30  31 32 33 34 35 36 37 38  |3456789012345678|
		000000e0  39 30 31 32 33 34 35 36  37 38 39 30 31 32 33 34  |9012345678901234|
		000000f0  35 36 37 38 39 30 31 32  33 34 35 36 37 38 39 30  |5678901234567890|
		00000100  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
		00000110  37 38 39 30 31 32 33 34  35 36 37 38 39 30 31 32  |7890123456789012|
		00000120  33 34 35 36 37 38 39                              |3456789|
	*/
	c.Convey("Given a message with missing option length extended value", t, func() {
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0D, // cut off
		}
		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then decode should fail with 'PacketIsTooShort' error", func() {
				c.So(err, c.ShouldEqual, PacketIsTooShort)
			})
		})
	})
}

func TestParseMessagePayload(t *testing.T) {
	c.Convey("Given a message with payload", t, func() {
		b := []byte{
			0x44, 0x02, 0x1B, 0x2B, 0x00, 0x00, 0x3F, 0x3D, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0B, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xFF, 0x61, 0x6C, 0x65, 0x78, 0x31, 0x32, 0x33,
		}
		c.Convey("When message is decoded", func() {
			m, _ := Decode(b, nil)
			c.Convey("Payload is as expected", func() {
				c.So(string(m.Payload), c.ShouldEqual, "alex123")
			})
		})
	})
}

func TestParseOptionWithLongExtendedLength(t *testing.T) {
	/*
		00000000  44 02 81 18 00 00 3a 31  39 6c 6f 63 61 6c 68 6f  |D.....:19localho|
		00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
		00000020  03 62 3d 55 06 6c 74 3d  33 30 30 0e 00 1f 30 31  |.b=U.lt=300...01|
		00000030  32 33 34 35 36 37 38 39  30 31 32 33 34 35 36 37  |2345678901234567|
		00000040  38 39 30 31 32 33 34 35  36 37 38 39 30 31 32 33  |8901234567890123|
		00000050  34 35 36 37 38 39 30 31  32 33 34 35 36 37 38 39  |4567890123456789|
		00000060  30 31 32 33 34 35 36 37  38 39 30 31 32 33 34 35  |0123456789012345|
		00000070  36 37 38 39 30 31 32 33  34 35 36 37 38 39 30 31  |6789012345678901|
		00000080  32 33 34 35 36 37 38 39  30 31 32 33 34 35 36 37  |2345678901234567|
		00000090  38 39 30 31 32 33 34 35  36 37 38 39 30 31 32 33  |8901234567890123|
		000000a0  34 35 36 37 38 39 30 31  32 33 34 35 36 37 38 39  |4567890123456789|
		000000b0  30 31 32 33 34 35 36 37  38 39 30 31 32 33 34 35  |0123456789012345|
		000000c0  36 37 38 39 30 31 32 33  34 35 36 37 38 39 30 31  |6789012345678901|
		000000d0  32 33 34 35 36 37 38 39  30 31 32 33 34 35 36 37  |2345678901234567|
		000000e0  38 39 30 31 32 33 34 35  36 37 38 39 30 31 32 33  |8901234567890123|
		000000f0  34 35 36 37 38 39 30 31  32 33 34 35 36 37 38 39  |4567890123456789|
		00000100  30 31 32 33 34 35 36 37  38 39 30 31 32 33 34 35  |0123456789012345|
		00000110  36 37 38 39 30 31 32 33  34 35 36 37 38 39 30 31  |6789012345678901|
		00000120  32 33 34 35 36 37 38 39  30 31 32 33 34 35 36 37  |2345678901234567|
		00000130  38 39 30 31 32 33 34 35  36 37 38 39 30 31 32 33  |8901234567890123|
		00000140  34 35 36 37 38 39 30 31  32 33 34 35 36 37 38 39  |4567890123456789|
		00000150  30 31 32 33 34 35 36 37  38 39                    |0123456789|
	*/
	c.Convey("Given a message with very long option", t, func() {
		b := []byte{
			0x44, 0x02, 0x81, 0x18, 0x00, 0x00, 0x3A, 0x31, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0E, 0x00, 0x1F, 0x30, 0x31,
			0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
			0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
			0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,
			0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
			0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
			0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,
			0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
			0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
			0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,
			0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
			0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		}
		c.Convey("When decoded", func() {
			m, err := Decode(b, nil)
			c.Convey("Then Uri-Query[3] should be as expected", func() {
				expected := strings.Join(
					[]string{
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
						"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
					}, "")

				c.So(err, c.ShouldBeNil)

				v, ok := (*m.Options)[UriQuery]
				c.So(ok, c.ShouldBeTrue)
				c.ShouldEqual(string(v[3]), expected)
			})
		})
	})
}

func TestMissingLongExtendedLengthBytesCauseError(t *testing.T) {
	c.Convey("Given a message with broken extended option length", t, func() {
		b := []byte{
			0x44, 0x02, 0x81, 0x18, 0x00, 0x00, 0x3A, 0x31, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0E, 0x00, // cut off on length
		}
		c.Convey("When decoded", func() {
			_, err := Decode(b, nil)
			c.Convey("Then the error should be 'PacketIsTooShort'", func() {
				c.ShouldEqual(err, PacketIsTooShort)
			})
		})
	})
}

func TestInvalidOptionNumber(t *testing.T) {
	c.Convey("Given a message with invalid option number", t, func() {
		/*
			00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
			00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
			00000020  03 62 3d 55 56 6c 74 3d  33 30 30                 |.b=U.lt=300     |
		*/
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x16, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30,
		}
		c.Convey("When message is decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then error should show 'InvalidOptionNumber'", func() {
				c.So(err, c.ShouldEqual, InvalidOptionNumber)
			})
		})
	})
}

func TestMessageFormatError(t *testing.T) {
	c.Convey("Given a message with invalid option delta", t, func() {
		/*
			00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
			00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
			00000020  03 62 3d 55 56 6c 74 3d  33 30 30                 |.b=U.lt=300     |
		*/
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0xF6, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30,
		}
		c.Convey("When message is decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then error should show 'MessageFormatError'", func() {
				c.So(err, c.ShouldEqual, MessageFormatError)
			})
		})
	})
}

func TestMessageFormatErrorOnInvalidOptionLength(t *testing.T) {
	c.Convey("Given a message with invalid option length", t, func() {
		/*
			00000000  44 02 ec 8e 00 00 e8 17  39 6c 6f 63 61 6c 68 6f  |D.......9localho|
			00000010  73 74 42 16 33 42 72 64  47 65 70 3d 61 6c 65 78  |stB.3BrdGep=alex|
			00000020  03 62 3d 55 56 6c 74 3d  33 30 30                 |.b=U.lt=300     |
		*/
		b := []byte{
			0x44, 0x02, 0x5D, 0x28, 0x00, 0x00, 0x82, 0x1C, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x5F, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30,
		}
		c.Convey("When message is decoded", func() {
			_, err := Decode(b, nil)

			c.Convey("Then error should show 'MessageFormatError'", func() {
				c.So(err, c.ShouldEqual, MessageFormatError)
			})
		})
	})
}

func TestEncodeDecodedMessageGivesTheSameByteContent(t *testing.T) {
	c.Convey("Given a message", t, func() {
		b := []byte{
			0x44, 0x02, 0x1B, 0x2B, 0x00, 0x00, 0x3F, 0x3D, 0x39, 0x6C, 0x6F, 0x63, 0x61, 0x6C, 0x68, 0x6F,
			0x73, 0x74, 0x42, 0x16, 0x33, 0x42, 0x72, 0x64, 0x47, 0x65, 0x70, 0x3D, 0x61, 0x6C, 0x65, 0x78,
			0x03, 0x62, 0x3D, 0x55, 0x06, 0x6C, 0x74, 0x3D, 0x33, 0x30, 0x30, 0x0B, 0x30, 0x31, 0x32, 0x33,
			0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xFF, 0x61, 0x6C, 0x65, 0x78, 0x31, 0x32, 0x33,
		}
		c.Convey("When decoded", func() {
			m, _ := Decode(b, nil)

			c.Convey("And encoded again", func() {
				b2 := m.Encode()
				c.Convey("Then the encoded message gives the same byte content", func() {
					c.So(b2, c.ShouldResemble, b)
				})
			})
		})
	})
}
