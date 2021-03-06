package coap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/aellwein/slf4go"
)

/* ERRORS */
var (
	PacketIsTooShort      = errors.New("packet is too short")
	InvalidMessageVersion = errors.New("invalid message version")
	InvalidTokenLength    = errors.New("invalid token length")
	MessageFormatError    = errors.New("message format error")
	InvalidOptionNumber   = errors.New("invalid option number")
)

/* MESSAGE */

// Message ID, 16 bits.
type MessageIdType uint16

func (m MessageIdType) String() string {
	return fmt.Sprintf("0x%04X", uint16(m))
}

func NewMessageId() MessageIdType {
	return MessageIdType(rand.Uint32())
}

// Version identifier. The only accepted version.
const MessageVersion byte = 0x01

// Basic message type, fits in 2 bits.
type MessageType int8

// Message types used in CoAP.
const (
	Confirmable MessageType = iota
	NonConfirmable
	Acknowledgement
	Reset
)

var AllMessageTypes = []MessageType{Confirmable, NonConfirmable, Acknowledgement, Reset}

type CodeClassType uint8
type CodeDetailType uint8

// Code: class, detail (c.dd)
type CodeType struct {
	CodeClass  CodeClassType
	CodeDetail CodeDetailType
}

// The CoAP message.
type Message struct {
	Type      MessageType
	Code      *CodeType
	MessageID MessageIdType
	Token     *TokenType
	Source    *net.UDPAddr
	Options   *OptionsType
	Payload   *PayloadType
}

// message type to string
func (mt MessageType) String() string {
	switch mt {
	case Confirmable:
		return "CON"
	case NonConfirmable:
		return "NON"
	case Acknowledgement:
		return "ACK"
	case Reset:
		return "RST"
	default:
		logger := slf4go.GetLogger("message")
		logger.Panicf("unknown message type: %d", mt)
		return ""
	}
}

/* TOKEN */

// Token, max 8 bytes.
type TokenType []byte

func (t TokenType) String() string {
	var b bytes.Buffer
	b.WriteString("0x")
	for _, i := range t {
		b.WriteString(fmt.Sprintf("%02X", i))
	}
	return b.String()
}

func NewToken() *TokenType {
	t := TokenType(make([]byte, 8))
	binary.BigEndian.PutUint64(t, rand.Uint64())
	return &t
}

func (t *TokenType) Copy() *TokenType {
	to := make([]byte, 8)
	copy(to, *t)
	return t
}

/* PAYLOAD */
type PayloadType struct {
	Type    *ContentType
	Content []byte
}

func (p *PayloadType) String() string {
	switch *p.Type {

	case ContentTypeTextPlain,
		ContentTypeApplicationLinkFormat,
		ContentTypeApplicationXml,
		ContentTypeApplicationJson:
		return string(p.Content)

	default:
		return HexContent(p.Content)
	}
}

// Initializes random number generator.
func init() {
	// Needed for Token/Message ID generation
	rand.Seed(time.Now().UnixNano())
}

// Reads and parses a CoAP Message from packet
func decode(buffer []byte, peer *net.UDPAddr) (*Message, error) {

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
		tkn = TokenType(make([]byte, tokenLength))
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
	var payload *PayloadType

	if payloadLen > 0 {
		payload = new(PayloadType)
		payload.Content = make([]byte, payloadLen-1)
		copy(payload.Content, buffer[pos+1:])

		if v, ok := opts[ContentFormat]; ok {
			var c ContentType
			switch len(v[0]) {
			case 0:
				c = 0
				payload.Type = &c
			default:
				if be, err := ToBigEndianNumber(v[0]); err == nil {
					if b, casts := be.(uint8); casts {
						c = ContentType(b)
					} else if u, casts := be.(uint16); casts {
						c = ContentType(u)
					} else {
						return nil, errors.New("invalid content format")
					}
					payload.Type = &c
				} else {
					return nil, err
				}
			}
		}
		// the case that the content format is not provided
		// is handled in message.Validate()
	}

	msg := &Message{
		Type: MessageType(mType),
		Code: &CodeType{
			CodeClass:  CodeClassType(codeClass),
			CodeDetail: CodeDetailType(codeDetail),
		},
		MessageID: MessageIdType(messageId),
		Token:     &tkn,
		Source:    peer,
		Options:   &opts,
		Payload:   payload,
	}

	return msg, nil
}

// Encode the message to a byte array.
func (m *Message) ToBytes() []byte {
	var pkt bytes.Buffer

	pkt.WriteByte(byte(64 + byte(m.Type<<4) + byte(len(*m.Token))))
	pkt.WriteByte(byte(m.Code.CodeClass<<5) + byte(m.Code.CodeDetail))

	msgId := make([]byte, 2)
	binary.BigEndian.PutUint16(msgId, uint16(m.MessageID))
	pkt.Write(msgId)
	pkt.Write(*m.Token)

	pkt.Write(encodeOptions(m.Options))

	if m.Payload != nil && len(m.Payload.Content) > 0 {
		pkt.WriteByte(0xff)
		pkt.Write(m.Payload.Content)
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

// Returns true, if message contains an option of given code.
func (m *Message) HasOption(opt OptionNumberType) bool {
	_, ok := (*m.Options)[opt]
	return ok
}

// Validates the message, returning one of the ok codes, if message is alright,
// otherwise specific error is returned.
func (m *Message) Validate() *CodeType {
	if m.Payload != nil && m.Payload.Content != nil && !m.HasOption(ContentFormat) {
		return BadRequest
	}
	return Ok
}
