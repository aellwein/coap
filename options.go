package coap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

/* OPTIONS */
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
				(*options)[optKey] = []OptionValueType{optionValue}
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
				if n, err := ToBigEndianNumber(i); err == nil {
					b.WriteString(fmt.Sprintf("%d", n))
					b.WriteString(",")
				}
			}
		}
		b.WriteString("] ")
	}
	return b.String()
}

/* OPTION NUMBERS */
type OptionFormat uint8

const (
	Empty OptionFormat = iota
	Opaque
	Uint
	String
)

type OptionDefinition struct {
	C       bool
	U       bool
	N       bool
	R       bool
	Name    string
	Format  OptionFormat
	Default interface{}
}

const (
	IfMatch       OptionNumberType = iota + 1
	UriHost                        = 3
	ETag                           = 4
	IfNoneMatch                    = 5
	UriPort                        = 7
	LocationPath                   = 8
	UriPath                        = 11
	ContentFormat                  = 12
	MaxAge                         = 14
	UriQuery                       = 15
	Accept                         = 17
	LocationQuery                  = 20
	ProxyUri                       = 35
	ProxyScheme                    = 39
	Size1                          = 60
)

// Lookup table for possible options.
var OptionLookupTable = map[OptionNumberType]OptionDefinition{
	IfMatch: {
		C:      true,
		U:      false,
		N:      false,
		R:      true,
		Name:   "If-Match",
		Format: Opaque,
	},
	UriHost: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Uri-Host",
		Format: String,
	},
	ETag: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "ETag",
		Format: Opaque,
	},
	IfNoneMatch: {
		C:      true,
		U:      false,
		N:      false,
		R:      false,
		Name:   "If-None-Match",
		Format: Empty,
	},
	UriPort: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Uri-Port",
		Format: Uint,
	},
	LocationPath: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "Location-Path",
		Format: String,
	},
	UriPath: {
		C:      true,
		U:      true,
		N:      false,
		R:      true,
		Name:   "Uri-Path",
		Format: String,
	},
	ContentFormat: {
		C:      false,
		U:      false,
		N:      false,
		R:      false,
		Name:   "Content-Format",
		Format: Uint,
	},
	MaxAge: {
		C:       false,
		U:       true,
		N:       false,
		R:       false,
		Name:    "Max-Age",
		Format:  Uint,
		Default: 60,
	},
	UriQuery: {
		C:      true,
		U:      true,
		N:      false,
		R:      true,
		Name:   "Uri-Query",
		Format: String,
	},
	Accept: {
		C:      true,
		U:      false,
		N:      false,
		R:      false,
		Name:   "Accept",
		Format: Uint,
	},
	LocationQuery: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "Location-Query",
		Format: Uint,
	},
	ProxyUri: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Proxy-Uri",
		Format: String,
	},

	ProxyScheme: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Proxy-Scheme",
		Format: Uint,
	},

	Size1: {
		C:      false,
		U:      false,
		N:      true,
		R:      false,
		Name:   "Size1",
		Format: Uint,
	},
}

// option number is in uint16 range
type OptionNumberType uint16

// NewLocationPathOption simplifies creation of a slash-style location path as a CoAP option.
func NewLocationPathOption(path string) []OptionValueType {
	pathElements := strings.Split(path, "/")
	ovt := make([]OptionValueType, 0)
	for _, elem := range pathElements {
		if elem != "" {
			ovt = append(ovt, OptionValueType(elem))
		}
	}
	return ovt
}

func UriPathOptionToString(opt []OptionValueType) string {
	var b bytes.Buffer
	for _, o := range opt {
		b.WriteString("/")
		b.WriteString(string(o))
	}
	return b.String()
}
