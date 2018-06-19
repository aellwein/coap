package message

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// option number is in uint16 range
type OptionNumberType uint16

// Option value consists of arbitrary bytes.
type OptionValueType []byte

// Options are build as map:
// OptionNumber : [ OptionValue1, OptionValue2, ...OptionValueN ]
type OptionsType map[OptionNumberType][]OptionValueType

// decode option from message buffer and return the rest of the buffer, if any
func decodeOptions(options OptionsType, buffer []byte) error {
	var (
		optionDelta  int
		optionLength int
		optionValue  []byte
	)
	if len(buffer) > 0 {
		i := 0
		for {
			od := buffer[i] >> 4
			ol := buffer[i] & 0xF

			// figure out Option Delta
			switch od {

			case 13:
				i++
				if len(buffer) < i+1 {
					return PacketIsTooShort
				}
				optionDelta += int(buffer[i])
				i++

			case 14:
				i++
				if len(buffer) < i+2 {
					return PacketIsTooShort
				}
				optionDelta = int(binary.BigEndian.Uint16(buffer[i:i+2]) + 269)
				i += 2

			case 15:
				// whole byte should be 0xFF, otherwise error
				if ol != 0xF {
					return MessageFormatError
				} else {
					// end of options detected
					return nil
				}
			default:
				optionDelta = int(od)

			}

			// figure out Option Length
			switch ol {

			case 13:
				i++
				if len(buffer) < i+1 {
					return PacketIsTooShort
				}
				optionLength += int(buffer[i])
				i++

			case 14:
				i++
				if len(buffer) < i+2 {
					return PacketIsTooShort
				}
				optionLength = int(binary.BigEndian.Uint16(buffer[i:i+2]) - 269)
				i += 2

			case 15:
				// invalid value
				return MessageFormatError

			default:
				optionLength = int(ol)
				i++
			}

			if len(buffer) < i+optionLength {
				return PacketIsTooShort
			}
			optionValue = make([]byte, optionLength)
			copy(optionValue, buffer[i:i+optionLength])
			fmt.Printf("optionDelta: %v, optionLength: %v, i: %v, optionValue: %s\n", optionDelta, optionLength, i, hex.Dump(optionValue))
		}
	}
	return nil
}
