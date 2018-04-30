package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Version identifier. The only accepted version.
const MessageVersion byte = 0x01

// The CoAP message.
type Message struct {
	Type        MessageType
	TokenLength TokenLengthType
	Code        CodeType
	MessageID   MessageIdType
	Token       TokenType
	Source      *net.UDPAddr
	Options     *OptionsType
}

var (
	PacketIsTooShort      = errors.New("packet is too short")
	InvalidMessageVersion = errors.New("invalid message version")
	InvalidTokenLength    = errors.New("invalid token length")
)

//func decodeOptions(buffer []byte, opts map[OptionNumber][]OptionValue) (error, bool, []byte) {
//
//}

// Reads and parses a CoAP Message from packet
func Decode(buffer []byte, peer *net.UDPAddr) (*Message, error) {

	if len(buffer) < 4 {
		// packet is too short
		return nil, PacketIsTooShort
	}

	version := buffer[0] >> 6

	if version != MessageVersion {
		// Spec: messages with unknown version numbers MUST be silently ignored.
		return nil, InvalidMessageVersion
	}

	mType := buffer[0] >> 4 & 3
	tokenLength := buffer[0] & 15

	if tokenLength > 8 {
		return nil, InvalidTokenLength
	}
	if len(buffer) < int(4+tokenLength) {
		// packet too short
		return nil, PacketIsTooShort
	}

	codeClass := buffer[1] >> 5
	codeDetail := buffer[1] & 31
	messageId := binary.BigEndian.Uint16(buffer[2:])

	tkn := uint64(0)

	if tokenLength != 0 {
		tkn = binary.BigEndian.Uint64(buffer[4 : 4+tokenLength])
	}

	// parse options, if any
	//if len(buffer) > int(4+tkl) {
	//	b := buffer[4+tkl:]
	//	opts := make(map[OptionNumber][]OptionValue)
	//	for {
	//		var (
	//			err error
	//			n   bool
	//		)
	//		err, n, b = decodeOptions(b, opts)
	//		if err != nil {
	//			return nil, err
	//		}
	//		if !n {
	//			break
	//		}
	//	}
	//}

	msg := &Message{
		Type:        MessageType(mType),
		TokenLength: TokenLengthType(tokenLength),
		Code: CodeType{
			CodeClass:  CodeClassType(codeClass),
			CodeDetail: CodeDetailType(codeDetail),
		},
		MessageID: MessageIdType(messageId),
		Token:     TokenType(tkn),
		Source:    peer,
	}

	return msg, nil
}

// Stringify message
func (m *Message) String() string {
	return fmt.Sprintf("Message{type=%v, code=%v, id=%v, tkn=0x%0X (%d), from=%v}",
		m.Type,
		m.Code,
		m.MessageID,
		m.Token,
		m.TokenLength,
		m.Source)
}
