package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/aellwein/coap/util"
)

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
			optionDelta = int(buffer[i]) + 13
			i++

		case 14:
			i++
			if len(buffer) < i+2 {
				return i, PacketIsTooShort
			}
			optionDelta = int(binary.BigEndian.Uint16(buffer[i:i+2]) + 269)
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
			optionLength = int(buffer[i]) + 13
			i++

		case 14:
			i++
			if len(buffer) < i+2 {
				return i, PacketIsTooShort
			}
			optionLength = int(binary.BigEndian.Uint16(buffer[i:i+2]) + 269)
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

		// check if option is valid
		// TODO: apply CUNR validation here.
		optKey := OptionNumberType(optionDelta)

		// key exists?
		if _, ok := OptionLookupTable[optKey]; ok {

			// options map already contains the option?
			if v, ok := (*options)[optKey]; ok {
				(*options)[optKey] = append(v, optionValue)
			} else {
				optValList := make([]OptionValueType, 1)
				optValList[0] = optionValue
				(*options)[optKey] = optValList
			}
		} else {
			return i, InvalidOptionNumber
		}
	}
	return i, nil
}

func (t OptionNumberType) String() string {
	return OptionLookupTable[t].Name
}

func (opt *OptionsType) String() string {
	var b bytes.Buffer
	for k, v := range *opt {
		b.WriteString(fmt.Sprintf("'%v'=[", k))
		for _, i := range v {
			switch OptionLookupTable[k].Format {

			case Empty:
				b.WriteString("{},")

			case Opaque:
				b.WriteString("[")
				b.WriteString(HexContent(i))
				b.WriteString("],")

			case String:
				b.WriteString("\"")
				b.WriteString(string(i))
				b.WriteString("\",")

			case Uint:
				n := util.ToBigEndianNumber(i)
				b.WriteString(fmt.Sprintf("%d", n))
				b.WriteString(",")
			}
		}
		b.WriteString("] ")
	}
	return b.String()
}
