package coap

import (
	"encoding/binary"
	"errors"
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

func ToBigEndianNumber(b []byte) (interface{}, error) {
	switch len(b) {
	case 0:
		return nil, errors.New("number of zero length is invalid")
	case 1:
		return b[0], nil
	case 2:
		return binary.BigEndian.Uint16(b), nil
	case 3:
		return binary.BigEndian.Uint32(padBE(b, 4)), nil
	case 4:
		return binary.BigEndian.Uint32(b), nil
	case 5, 6, 7:
		return binary.BigEndian.Uint64(padBE(b, 8)), nil
	case 8:
		return binary.BigEndian.Uint64(b), nil
	default:
		// everything above 8 bytes gets boxed into 8 bytes
		return binary.BigEndian.Uint64(b[len(b)-8:]), nil
	}
}

func ToLittleEndianNumber(b []byte) (interface{}, error) {
	switch len(b) {
	case 0:
		return nil, errors.New("number of zero length is invalid")
	case 1:
		return b[0], nil
	case 2:
		return binary.LittleEndian.Uint16(b), nil
	case 3:
		return binary.LittleEndian.Uint32(padLE(b, 4)), nil
	case 4:
		return binary.LittleEndian.Uint32(b), nil
	case 5, 6, 7:
		return binary.LittleEndian.Uint64(padLE(b, 8)), nil
	case 8:
		return binary.LittleEndian.Uint64(b), nil
	default:
		// everything above 8 bytes gets boxed into 8 bytes
		return binary.LittleEndian.Uint64(b[0:8]), nil
	}
}
