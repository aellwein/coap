package util

import (
	"encoding/binary"
)

func padBE(b []byte, targetLen int) []byte {
	target := make([]byte, targetLen)
	copy(target[targetLen-len(b):], b)
	return target
}

func padLE(b []byte, targetLen int) []byte {
	target := make([]byte, targetLen)
	copy(target, b)
	return target
}

func ToBigEndianNumber(b []byte) interface{} {
	switch len(b) {
	case 1:
		return b[0]
	case 2:
		return binary.BigEndian.Uint16(b)
	case 3:
		return binary.BigEndian.Uint32(padBE(b, 4))
	case 4:
		return binary.BigEndian.Uint32(b)
	case 5, 6, 7:
		return binary.BigEndian.Uint64(padBE(b, 8))
	case 8:
		return binary.BigEndian.Uint64(b)
	default:
		// everything above 8 bytes gets boxed into 8 bytes
		return binary.BigEndian.Uint64(b[len(b)-8:])
	}
}

func ToLittleEndianNumber(b []byte) interface{} {
	switch len(b) {
	case 1:
		return b[0]
	case 2:
		return binary.LittleEndian.Uint16(b)
	case 3:
		return binary.LittleEndian.Uint32(padLE(b, 4))
	case 4:
		return binary.LittleEndian.Uint32(b)
	case 5, 6, 7:
		return binary.LittleEndian.Uint64(padLE(b, 8))
	case 8:
		return binary.LittleEndian.Uint64(b)
	default:
		// everything above 8 bytes gets boxed into 8 bytes
		return binary.LittleEndian.Uint64(b[0:8])
	}
}
