package coap

import (
	"encoding/binary"
	"net"
)

type messageBuilder struct {
	msgCtx *messageContext
}

type messageCodeBuilder struct {
	msgCtx *messageContext
}

type messageIdBuilder struct {
	msgCtx *messageContext
}

type messageTokenBuilder struct {
	msgCtx *messageContext
}

type messageOptionBuilder struct {
	msgCtx *messageContext
}
type messagePayloadBuilder struct {
	msgCtx *messageContext
}

type messageContext struct {
	mType     MessageType
	code      *CodeType
	messageId MessageIdType
	token     *TokenType
	source    *net.UDPAddr
	options   *OptionsType
	payload   PayloadType
}

// NewConfirmableMessageBuilder creates a new message builder for a confirmable message.
func NewConfirmableMessageBuilder() messageBuilder {
	return messageBuilder{&messageContext{mType: Confirmable}}
}

// NewAcknowledgementMessageBuilder creates a new message builder for an acknowledgement message.
func NewAcknowledgementMessageBuilder() messageBuilder {
	return messageBuilder{&messageContext{mType: Acknowledgement}}
}

// NewResetMessageBuilder creates a new message builder for a reset message.
func NewResetMessageBuilder() messageBuilder {
	return messageBuilder{&messageContext{mType: Reset}}
}

// NewNonConfirmableMessageBuilder creates a new message builder for a non-confirmable message.
func NewNonConfirmableMessageBuilder() messageBuilder {
	return messageBuilder{&messageContext{mType: Reset}}
}

// Code builder method provides the message builder with message code.
func (m messageBuilder) Code(code *CodeType) messageCodeBuilder {
	m.msgCtx.code = code
	return messageCodeBuilder{m.msgCtx}
}

// From builder method (optional) provides the source address of where the message came from.
func (m messageBuilder) From(addr *net.UDPAddr) messageBuilder {
	m.msgCtx.source = addr
	return m
}

// MessageId builder method provides the message ID of the CoAP message.
func (m messageCodeBuilder) MessageId(messageId MessageIdType) messageIdBuilder {
	m.msgCtx.messageId = messageId
	return messageIdBuilder{m.msgCtx}
}

// WithRandomMessageId builder method sets a new, random messageID into the message.
func (m messageCodeBuilder) WithRandomMessageId() messageIdBuilder {
	m.msgCtx.messageId = NewMessageId()
	return messageIdBuilder{m.msgCtx}
}

// Token builder method sets the given token into the message.
func (m messageIdBuilder) Token(token *TokenType) messageTokenBuilder {
	m.msgCtx.token = token
	return messageTokenBuilder{m.msgCtx}
}

// WithRandomToken builder method sets a new, random token into the message.
func (m messageIdBuilder) WithRandomToken() messageTokenBuilder {
	m.msgCtx.token = NewToken()
	return messageTokenBuilder{m.msgCtx}
}

func (m messageTokenBuilder) Build() *Message {
	opts := make(OptionsType)
	return &Message{
		Type:      m.msgCtx.mType,
		Code:      m.msgCtx.code,
		MessageID: m.msgCtx.messageId,
		Token:     m.msgCtx.token,
		Source:    m.msgCtx.source,
		Options:   &opts,
	}
}

// Option builder method adds an option (code), with one or many values given.
func (m messageTokenBuilder) Option(opt OptionNumberType, valueTypes ...OptionValueType) messageOptionBuilder {
	if m.msgCtx.options == nil {
		opts := make(OptionsType)
		m.msgCtx.options = &opts
		(*m.msgCtx.options)[opt] = valueTypes
	} else {
		(*m.msgCtx.options)[opt] = append((*m.msgCtx.options)[opt], valueTypes...)
	}
	return messageOptionBuilder{m.msgCtx}
}

// Option builder method adds an option (code), with one or many values given.
func (m messageOptionBuilder) Option(opt OptionNumberType, valueTypes ...OptionValueType) messageOptionBuilder {
	if m.msgCtx.options == nil {
		opts := make(OptionsType)
		m.msgCtx.options = &opts
		(*m.msgCtx.options)[opt] = valueTypes
	} else {
		(*m.msgCtx.options)[opt] = append((*m.msgCtx.options)[opt], valueTypes...)
	}
	return messageOptionBuilder{m.msgCtx}
}

// Build constructs a new message.
func (m messageOptionBuilder) Build() *Message {
	return &Message{
		Type:      m.msgCtx.mType,
		Code:      m.msgCtx.code,
		MessageID: m.msgCtx.messageId,
		Token:     m.msgCtx.token,
		Source:    m.msgCtx.source,
		Options:   m.msgCtx.options,
	}
}

// WithPayload provides a payload of given type to the message builder.
func (m messageOptionBuilder) WithPayload(cType ContentType, payload PayloadType) messagePayloadBuilder {
	if m.msgCtx.options == nil {
		opts := make(OptionsType)
		m.msgCtx.options = &opts
	}
	if cType < 256 {
		(*m.msgCtx.options)[ContentFormat] = []OptionValueType{[]byte{byte(cType)}}
	} else {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(cType))
		(*m.msgCtx.options)[ContentFormat] = []OptionValueType{b}
	}
	m.msgCtx.payload = payload
	return messagePayloadBuilder{m.msgCtx}
}

// Build constructs a new message.
func (m messagePayloadBuilder) Build() *Message {
	return &Message{
		Type:      m.msgCtx.mType,
		Code:      m.msgCtx.code,
		MessageID: m.msgCtx.messageId,
		Token:     m.msgCtx.token,
		Source:    m.msgCtx.source,
		Options:   m.msgCtx.options,
		Payload:   m.msgCtx.payload,
	}
}

// NewMessageFromBytes constructs a new message from the given bytes packet,
// if not successful, an error is returned.
func NewMessageFromBytes(buffer []byte) (*Message, error) {
	return decode(buffer, nil)
}

// NewMessageFromBytesAndPeer constructs a new message from the given bytes packet and
// sets the source address of the message. If not successful, an error is returned.
func NewMessageFromBytesAndPeer(buffer []byte, peer *net.UDPAddr) (*Message, error) {
	return decode(buffer, peer)
}

func NewInternalServerErrorResponseMessage(msg *Message) *Message {
	return NewAcknowledgementMessageBuilder().
		Code(InternalServerError).
		MessageId(msg.MessageID).
		Token(msg.Token).
		Build()
}

func NewMethodNotAllowedResponseMessage(msg *Message) *Message {
	return NewAcknowledgementMessageBuilder().
		Code(MethodNotAllowed).
		MessageId(msg.MessageID).
		Token(msg.Token).
		Build()
}

func NewBadRequestResponseMessage(msg *Message) *Message {
	return NewAcknowledgementMessageBuilder().
		Code(BadRequest).
		MessageId(msg.MessageID).
		Token(msg.Token).
		Build()
}
func NewNotFoundResponseMessage(msg *Message) *Message {
	return NewAcknowledgementMessageBuilder().
		Code(NotFound).
		MessageId(msg.MessageID).
		Token(msg.Token).
		Build()
}
