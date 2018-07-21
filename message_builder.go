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
	payload   *PayloadType
}

// NewMessageBuilderOfType creates a new message builder for message of given type.
func NewMessageBuilderOfType(mt MessageType) messageBuilder {
	opts := make(OptionsType)
	return messageBuilder{&messageContext{mType: mt, options: &opts}}
}

// NewConfirmableMessageBuilder creates a new message builder for a confirmable message.
func NewConfirmableMessageBuilder() messageBuilder {
	return NewMessageBuilderOfType(Confirmable)
}

// NewAcknowledgementMessageBuilder creates a new message builder for an acknowledgement message.
func NewAcknowledgementMessageBuilder() messageBuilder {
	return NewMessageBuilderOfType(Acknowledgement)
}

// NewResetMessageBuilder creates a new message builder for a reset message.
func NewResetMessageBuilder() messageBuilder {
	return NewMessageBuilderOfType(Reset)
}

// NewNonConfirmableMessageBuilder creates a new message builder for a non-confirmable message.
func NewNonConfirmableMessageBuilder() messageBuilder {
	return NewMessageBuilderOfType(NonConfirmable)
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
	if m.msgCtx.options == nil {
		opts := make(OptionsType)
		m.msgCtx.options = &opts
	}
	return &Message{
		Type:      m.msgCtx.mType,
		Code:      m.msgCtx.code,
		MessageID: m.msgCtx.messageId,
		Token:     m.msgCtx.token,
		Source:    m.msgCtx.source,
		Options:   m.msgCtx.options,
	}
}

// Option builder method adds an option (code), with one or many values given.
func (m messageTokenBuilder) Option(opt OptionNumberType, valueTypes ...OptionValueType) messageTokenBuilder {
	if m.msgCtx.options == nil {
		opts := make(OptionsType)
		m.msgCtx.options = &opts
		(*m.msgCtx.options)[opt] = valueTypes
	} else {
		(*m.msgCtx.options)[opt] = append((*m.msgCtx.options)[opt], valueTypes...)
	}
	return m
}

// WithPayload provides a payload of given type to the message builder.
func (m messageTokenBuilder) WithPayload(cType ContentType, payload []byte) messagePayloadBuilder {
	if cType < 256 {
		(*m.msgCtx.options)[ContentFormat] = []OptionValueType{[]byte{byte(cType)}}
	} else {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(cType))
		(*m.msgCtx.options)[ContentFormat] = []OptionValueType{b}
	}
	m.msgCtx.payload = &PayloadType{
		Type:    &cType,
		Content: payload,
	}
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

func responseWithCode(request *Message, code *CodeType) *Message {
	return NewAcknowledgementMessageBuilder().Code(code).MessageId(request.MessageID).Token(request.Token).Build()
}

// Content Response
func NewContentResponseMessage(request *Message) *Message {
	return responseWithCode(request, Content)
}

// Bad Request Response
func NewBadRequestResponseMessage(request *Message) *Message {
	return responseWithCode(request, BadRequest)
}

// Unauthorized Response
func NewUnauthorizedResponseMessage(request *Message) *Message {
	return responseWithCode(request, Unauthorized)
}

// Bad Option Response
func NewBadOptionResponseMessage(request *Message) *Message {
	return responseWithCode(request, BadOption)
}

// Forbidden Response
func NewForbiddenResponseMessage(request *Message) *Message {
	return responseWithCode(request, Forbidden)
}

// Not Found Response
func NewNotFoundResponseMessage(request *Message) *Message {
	return responseWithCode(request, NotFound)
}

// Method Not Allowed Response
func NewMethodNotAllowedResponseMessage(request *Message) *Message {
	return responseWithCode(request, MethodNotAllowed)
}

// Not Acceptable Response
func NewNotAcceptableResponseMessage(request *Message) *Message {
	return responseWithCode(request, NotAcceptable)
}

// Precondition Failed Response
func NewPreconditionFailedResponseMessage(request *Message) *Message {
	return responseWithCode(request, PreconditionFailed)
}

// Request Entity Too Large Response
func NewRequestEntityTooLargeResponseMessage(request *Message) *Message {
	return responseWithCode(request, RequestEntityTooLarge)
}

// Unsupported Content Format Response
func NewUnsupportedContentFormatResponseMessage(request *Message) *Message {
	return responseWithCode(request, UnsupportedContentFormat)
}

// Internal Server Error Response
func NewInternalServerErrorResponseMessage(request *Message) *Message {
	return responseWithCode(request, InternalServerError)
}

// Not Implemented Response
func NewNotImplementedResponseMessage(request *Message) *Message {
	return responseWithCode(request, NotImplemented)
}

// Bad Gateway Response
func NewBadGatewayResponseMessage(request *Message) *Message {
	return responseWithCode(request, BadGateway)
}

// Service Unavailable Response
func NewServiceUnavailableResponseMessage(request *Message) *Message {
	return responseWithCode(request, ServiceUnavailable)
}

// Gateway Timeout Response
func NewGatewayTimeoutResponseMessage(request *Message) *Message {
	return responseWithCode(request, GatewayTimeout)
}

// Proxying Not Supported Response
func NewProxyingNotSupportedResponseMessage(request *Message) *Message {
	return responseWithCode(request, ProxyingNotSupported)
}
