package message

import (
	"encoding/binary"
)

// option number is in uint16 range
type OptionNumberType uint16

// Option value consists of arbitrary bytes.
type OptionValueType []byte

// Options are build as map:
// OptionNumber : [ OptionValue1, OptionValue2, ...OptionValueN ]
type OptionsType map[OptionNumberType][]OptionValueType

// decode option from message buffer and return the next position in the buffer
func decodeOptions(options *OptionsType, buffer []byte) (int, error) {
	var (
		optionDelta  int
		optionLength int
		optionValue  []byte
	)
	i := 0
	for len(buffer) > i {
		od := buffer[i] >> 4
		ol := buffer[i] & 0xF

		// figure out Option Delta
		switch od {
		case 13:
			i++
			if len(buffer) < i+1 {
				return i, PacketIsTooShort
			}
			optionDelta += int(buffer[i])
			i++

		case 14:
			i++
			if len(buffer) < i+2 {
				return i, PacketIsTooShort
			}
			optionDelta += int(binary.BigEndian.Uint16(buffer[i:i+2]) + 269)
			i += 2

		case 15:
			// whole byte should be 0xFF, otherwise error
			if ol != 0xF {
				return i, MessageFormatError
			} else {
				// end of options detected
				return i, nil
			}
		default:
			optionDelta += int(od)
		}

		// figure out Option Length
		switch ol {

		case 13:
			i++
			if len(buffer) < i+1 {
				return i, PacketIsTooShort
			}
			optionLength += int(buffer[i])
			i++

		case 14:
			i++
			if len(buffer) < i+2 {
				return i, PacketIsTooShort
			}
			optionLength = int(binary.BigEndian.Uint16(buffer[i:i+2]) - 269)
			i += 2

		case 15:
			// invalid value
			return i, MessageFormatError

		default:
			optionLength = int(ol)
			i++
		}

		if len(buffer) < i+optionLength {
			return i, PacketIsTooShort
		}
		optionValue = make([]byte, optionLength)
		copy(optionValue, buffer[i:i+optionLength])
		i += optionLength
	}
	return i, nil
}
