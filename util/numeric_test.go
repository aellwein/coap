package util

import (
	"encoding/binary"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestToBigEndianNumberOfOneByte(t *testing.T) {
	Convey("Given a byte array with 1 element (0x20)", t, func() {
		var b = []byte{0x20}
		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)
			Convey("Result must be 0x20", func() {
				So(n.(byte), ShouldEqual, 0x20)
			})
		})
	})
}

func TestToLittleEndianNumberOfOneByte(t *testing.T) {
	Convey("Given a byte array with 1 element (0x20)", t, func() {
		var b = []byte{0x20}
		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)
			Convey("Result must be 0x20", func() {
				So(n.(byte), ShouldEqual, 0x20)
			})
		})
	})
}

func TestToBigEndianNumberOfTwoBytes(t *testing.T) {
	Convey("Given a big endian number of 2 bytes", t, func() {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, 5683)

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint16), ShouldEqual, 5683)
			})
		})
	})
}

func TestToLittleEndianNumberOfTwoBytes(t *testing.T) {
	Convey("Given a little endian number of 2 bytes", t, func() {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, 5683)

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint16), ShouldEqual, 5683)
			})
		})
	})
}

func TestToBigEndianNumberOfThreeBytes(t *testing.T) {
	Convey("Given a big endian number of 3 bytes", t, func() {
		b := []byte{0x00, 0x16, 0x33}

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint32), ShouldEqual, 5683)
			})
		})
	})
}

func TestToLittleEndianNumberOfThreeBytes(t *testing.T) {
	Convey("Given a little endian number of 3 bytes", t, func() {
		b := []byte{0x33, 0x16, 0x00}

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint32), ShouldEqual, 5683)
			})
		})
	})
}

func TestToBigEndianNumberOfFourBytes(t *testing.T) {
	Convey("Given a big endian number of 4 bytes", t, func() {
		b := []byte{0xCA, 0xFE, 0xBA, 0xBE}

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint32), ShouldEqual, 0xCAFEBABE)
			})
		})
	})
}

func TestToLittleEndianNumberOfFourBytes(t *testing.T) {
	Convey("Given a little endian number of 4 bytes", t, func() {
		b := []byte{0xBE, 0xBA, 0xFE, 0xCA}

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint32), ShouldEqual, 0xCAFEBABE)
			})
		})
	})
}

func TestToBigEndianNumberOfSixBytes(t *testing.T) {
	Convey("Given a big endian number of 6 bytes", t, func() {
		b := []byte{0x00, 0x00, 0xCA, 0xFE, 0xBA, 0xBE}

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint64), ShouldEqual, 0xCAFEBABE)
			})
		})
	})
}

func TestToLittleEndianNumberOfSixBytes(t *testing.T) {
	Convey("Given a little endian number of 6 bytes", t, func() {
		b := []byte{0xBE, 0xBA, 0xFE, 0xCA, 0x00, 0x00}

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint64), ShouldEqual, 0xCAFEBABE)
			})
		})
	})
}

func TestToBigEndianNumberOfEightBytes(t *testing.T) {
	Convey("Given a big endian number of 8 bytes", t, func() {
		b := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF}

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint64), ShouldEqual, uint64(0xCAFEBABEDEADBEEF))
			})
		})
	})
}

func TestToLittleEndianNumberOfEightBytes(t *testing.T) {
	Convey("Given a little endian number of 8 bytes", t, func() {
		b := []byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA}

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint64), ShouldEqual, uint64(0xCAFEBABEDEADBEEF))
			})
		})
	})
}

func TestToBigEndianNumberOfTenBytes(t *testing.T) {
	Convey("Given a big endian number of 10 bytes", t, func() {
		b := []byte{0x00, 0x00, 0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF}

		Convey("When converted to big endian number", func() {
			n := ToBigEndianNumber(b)

			Convey("Then the result must be equal the boxed uint64 value", func() {
				So(n.(uint64), ShouldEqual, uint64(0xCAFEBABEDEADBEEF))
			})
		})
	})
}

func TestToLittleEndianNumberOfTenBytes(t *testing.T) {
	Convey("Given a little endian number of 10 bytes", t, func() {
		b := []byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA, 0x00, 0x00}

		Convey("When converted to little endian number", func() {
			n := ToLittleEndianNumber(b)

			Convey("Then the result must be equal the given number", func() {
				So(n.(uint64), ShouldEqual, uint64(0xCAFEBABEDEADBEEF))
			})
		})
	})
}
