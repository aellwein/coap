package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/aellwein/coap/util"
)

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
