package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/aellwein/coap/util"
	"sort"
)

// Option value consists of arbitrary bytes.
type OptionValueType []byte

// Options are build as map:
// OptionNumber : [ OptionValue1, OptionValue2, ...OptionValueN ]
type OptionsType map[OptionNumberType][]OptionValueType

var EmptyOptions = make(OptionsType)

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
				// Spec: "The presence of a marker followed by a zero-length payload MUST be processed as a
				// message format error."
				if len(buffer) == i+1 {
					return i, MessageFormatError
				}
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

func encodeNum(num uint16) (byte, []byte) {
	switch {
	case num < 13:
		return byte(num), []byte{}
	case num < 269:
		return 13, []byte{byte(num - 13)}
	default:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(num-269))
		return 14, b
	}
}

func getOptionLengthBytes(opt OptionValueType) (byte, []byte) {
	return encodeNum(uint16(len(opt)))
}

func getOptionDeltaBytes(opt OptionNumberType, prev OptionNumberType) (byte, []byte) {
	optDelta := opt - prev
	return encodeNum(uint16(optDelta))
}

func encodeOptions(options *OptionsType) []byte {
	var b bytes.Buffer

	if len(*options) > 0 {

		// to build proper delta, we sort the option numbers (map keys)
		var optionNumbers []int
		for k := range *options {
			optionNumbers = append(optionNumbers, int(k))
		}
		sort.Ints(optionNumbers)

		var optDeltaPrev OptionNumberType = 0
		for _, i := range optionNumbers {
			for _, j := range (*options)[OptionNumberType(i)] {
				optDeltaByte, optDeltaBytes := getOptionDeltaBytes(OptionNumberType(i), optDeltaPrev)
				optDeltaPrev = OptionNumberType(i)
				optLenByte, optLenBytes := getOptionLengthBytes(j)
				b.WriteByte((optDeltaByte << 4) + optLenByte)
				if len(optDeltaBytes) > 0 {
					b.Write(optDeltaBytes)
				}
				if len(optLenBytes) > 0 {
					b.Write(optLenBytes)
				}
				b.Write(j)
			}
		}
	}
	return b.Bytes()
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
