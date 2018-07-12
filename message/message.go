package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Version identifier. The only accepted version.
const MessageVersion byte = 0x01

// The CoAP message.
type Message struct {
	Type      MessageType
	Code      CodeType
	MessageID MessageIdType
	Token     TokenType
	Source    *net.UDPAddr
	Options   *OptionsType
	Payload   PayloadType
}

var (
	PacketIsTooShort      = errors.New("packet is too short")
	InvalidMessageVersion = errors.New("invalid message version")
	InvalidTokenLength    = errors.New("invalid token length")
	MessageFormatError    = errors.New("message format error")
	InvalidOptionNumber   = errors.New("invalid option number")
)

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

	var tkn TokenType
	if tokenLength != 0 {
		tkn = make([]byte, tokenLength)
		copy(tkn, buffer[4:4+tokenLength])
	} else {
		tkn = []byte{}
	}

	opts := make(OptionsType)
	buf := buffer[4+int(tokenLength):]

	// parse options, if any
	pos, err := decodeOptions(&opts, buf)
	if err != nil {
		return nil, err
	}

	// parse payload, if any
	pos += int(tokenLength) + 4
	payloadLen := len(buffer) - pos
	var payload PayloadType
	if payloadLen > 0 {
		payload = make([]byte, payloadLen-1)
		copy(payload, buffer[pos+1:])
	}

	msg := &Message{
		Type: MessageType(mType),
		Code: CodeType{
			CodeClass:  CodeClassType(codeClass),
			CodeDetail: CodeDetailType(codeDetail),
		},
		MessageID: MessageIdType(messageId),
		Token:     tkn,
		Source:    peer,
		Options:   &opts,
		Payload:   payload,
	}

	return msg, nil
}

// Encode the message to a byte array.
func (m *Message) Encode() []byte {
	var pkt bytes.Buffer

	pkt.WriteByte(byte(64 + byte(m.Type<<4) + byte(len(m.Token))))
	pkt.WriteByte(byte(m.Code.CodeClass<<5) + byte(m.Code.CodeDetail))

	msgId := make([]byte, 2)
	binary.BigEndian.PutUint16(msgId, uint16(m.MessageID))
	pkt.Write(msgId)
	pkt.Write(m.Token)

	pkt.Write(encodeOptions(m.Options))

	if len(m.Payload) > 0 {
		pkt.WriteByte(0xff)
		pkt.Write(m.Payload)
	}
	return pkt.Bytes()
}

// Stringify message
func (m *Message) String() string {
	return fmt.Sprintf("Message{type=%v, code=%v, id=%v, tkn=%v, options=%v, payload=%v, from=%v}",
		m.Type,
		m.Code,
		m.MessageID,
		m.Token,
		m.Options,
		m.Payload,
		m.Source)
}

func HexContent(p []byte) string {
	var b bytes.Buffer
	b.WriteString("[")
	for n, i := range p {
		if n == 0 {
			b.WriteString(fmt.Sprintf("0x%02X", i))
		} else {
			b.WriteString(fmt.Sprintf(", 0x%02X", i))
		}
	}
	b.WriteString("]")
	return b.String()
}
