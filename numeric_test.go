package coap

import (
	"fmt"
	"reflect"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

type testMatrix struct {
	input        []byte
	Type         reflect.Type
	bigEndian    interface{}
	littleEndian interface{}
}

var (
	tByte   byte
	tUint16 uint16
	tUint32 uint32
	tUint64 uint64

	tests = []testMatrix{
		{
			input:        []byte{0x20},
			Type:         reflect.TypeOf(tByte),
			bigEndian:    0x20,
			littleEndian: 0x20,
		},
		{
			input:        []byte{0x13, 0x37},
			Type:         reflect.TypeOf(tUint16),
			bigEndian:    0x1337,
			littleEndian: 0x3713,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA},
			Type:         reflect.TypeOf(tUint32),
			bigEndian:    0xCAFEBA,
			littleEndian: 0x00BAFECA,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE},
			Type:         reflect.TypeOf(tUint32),
			bigEndian:    0xCAFEBABE,
			littleEndian: 0xBEBAFECA,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE},
			Type:         reflect.TypeOf(tUint64),
			bigEndian:    0xCAFEBABEDE,
			littleEndian: 0xDEBEBAFECA,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD},
			Type:         reflect.TypeOf(tUint64),
			bigEndian:    0xCAFEBABEDEAD,
			littleEndian: 0xADDEBEBAFECA,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE},
			Type:         reflect.TypeOf(tUint64),
			bigEndian:    0xCAFEBABEDEADBE,
			littleEndian: 0xBEADDEBEBAFECA,
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF},
			Type:         reflect.TypeOf(tUint64),
			bigEndian:    uint64(0xCAFEBABEDEADBEEF),
			littleEndian: uint64(0xEFBEADDEBEBAFECA),
		},
		{
			input:        []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF, 0x13, 0x37},
			Type:         reflect.TypeOf(tUint64),
			bigEndian:    uint64(0xBABEDEADBEEF1337),
			littleEndian: uint64(0xEFBEADDEBEBAFECA),
		},
	}
)

func TestToNumber_All(t *testing.T) {
	c.Convey("Given a byte array representing number", t, func() {
		for _, test := range tests {

			c.Convey(fmt.Sprintf("When converting %v to big endian number", HexContent(test.input)), func() {

				c.Convey(fmt.Sprintf("And array consists of %d bytes", len(test.input)), func() {

					n, _ := ToBigEndianNumber(test.input)

					c.Convey(fmt.Sprintf("Then resulting type is '%v'", test.Type), func() {
						c.So(reflect.TypeOf(n), c.ShouldEqual, test.Type)

						c.Convey("And the result equals the expected number", func() {
							c.So(n, c.ShouldEqual, test.bigEndian)
						})
					})
				})
			})
			c.Convey(fmt.Sprintf("When converting %v to little endian number", HexContent(test.input)), func() {
				c.Convey(fmt.Sprintf("And array consists of %d bytes", len(test.input)), func() {

					n, _ := ToLittleEndianNumber(test.input)

					c.Convey(fmt.Sprintf("Then resulting type is '%v'", test.Type), func() {
						c.So(reflect.TypeOf(n), c.ShouldEqual, test.Type)

						c.Convey("And the result equals the expected number", func() {
							c.So(n, c.ShouldEqual, test.littleEndian)
						})
					})
				})

			})
		}

		// special cases

		c.Convey("When a byte array of zero length is converted to big endian number", func() {
			_, err := ToBigEndianNumber([]byte{})
			c.Convey("Then an error is returned", func() {
				c.So(err, c.ShouldNotBeNil)
			})
		})

		c.Convey("When a byte array of zero length is converted to little endian number", func() {
			_, err := ToLittleEndianNumber([]byte{})
			c.Convey("Then an error is returned", func() {
				c.So(err, c.ShouldNotBeNil)
			})
		})
	})
}
