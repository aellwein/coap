package message

import "fmt"

// MessageCode
type MessageCode struct {
	cat  uint8
	code uint8
}

// message codes as defined by IANA
var (
	EmptyMessage             *MessageCode = &MessageCode{0, 00}
	GET                      *MessageCode = &MessageCode{0, 01}
	POST                     *MessageCode = &MessageCode{0, 02}
	PUT                      *MessageCode = &MessageCode{0, 03}
	DELETE                   *MessageCode = &MessageCode{0, 04}
	Created                  *MessageCode = &MessageCode{2, 01}
	Deleted                  *MessageCode = &MessageCode{2, 02}
	Valid                    *MessageCode = &MessageCode{2, 03}
	Changed                  *MessageCode = &MessageCode{2, 04}
	Content                  *MessageCode = &MessageCode{2, 05}
	BadRequest               *MessageCode = &MessageCode{4, 00}
	Unauthorized             *MessageCode = &MessageCode{4, 01}
	BadOption                *MessageCode = &MessageCode{4, 02}
	Forbidden                *MessageCode = &MessageCode{4, 03}
	NotFound                 *MessageCode = &MessageCode{4, 04}
	MethodNotAllowed         *MessageCode = &MessageCode{4, 05}
	NotAcceptable            *MessageCode = &MessageCode{4, 06}
	PreconditionFailed       *MessageCode = &MessageCode{4, 12}
	RequestEntityTooLarge    *MessageCode = &MessageCode{4, 13}
	UnsupportedContentFormat *MessageCode = &MessageCode{4, 15}
	InternalServerError      *MessageCode = &MessageCode{5, 00}
	NotImplemented           *MessageCode = &MessageCode{5, 01}
	BadGateway               *MessageCode = &MessageCode{5, 02}
	ServiceUnavailable       *MessageCode = &MessageCode{5, 03}
	GatewayTimeout           *MessageCode = &MessageCode{5, 04}
	ProxyingNotSupported     *MessageCode = &MessageCode{5, 05}
)

func (m *MessageCode) IsEmpty() bool {
	return m.cat == 0 && m.code == 0
}

func (m *MessageCode) IsRequest() bool {
	return m.cat == 0 && (m.code >= 1 && m.code <= 31)
}

func (m *MessageCode) IsResponse() bool {
	return (m.cat >= 2 && m.cat <= 5) && (m.code >= 0 && m.code <= 31)
}

func (m *MessageCode) String() string {
	return fmt.Sprintf("%d.%02d", m.cat, m.code)
}
