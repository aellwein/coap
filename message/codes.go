package message

import (
	"fmt"
	"github.com/aellwein/slf4go"
)

type CodeClassType uint8
type CodeDetailType uint8

// Code: class, detail (c.dd)
type CodeType struct {
	CodeClass  CodeClassType
	CodeDetail CodeDetailType
}

// Stringify code
func (c CodeType) String() string {
	return fmt.Sprintf("%d.%02d", c.CodeClass, c.CodeDetail)
}

type MethodCode CodeType
type ResponseCode CodeType

// Empty message code
var EmptyMessage = &CodeType{CodeClass: 0, CodeDetail: 0}

// Method Codes
var (
	GET    = &MethodCode{CodeClass: 0, CodeDetail: 1}
	POST   = &MethodCode{CodeClass: 0, CodeDetail: 2}
	PUT    = &MethodCode{CodeClass: 0, CodeDetail: 3}
	DELETE = &MethodCode{CodeClass: 0, CodeDetail: 4}
)

var (
	// success codes
	Ok      = &ResponseCode{CodeClass: 2, CodeDetail: 0}
	Created = &ResponseCode{CodeClass: 2, CodeDetail: 1}
	Deleted = &ResponseCode{CodeClass: 2, CodeDetail: 2}
	Valid   = &ResponseCode{CodeClass: 2, CodeDetail: 3}
	Changed = &ResponseCode{CodeClass: 2, CodeDetail: 4}
	Content = &ResponseCode{CodeClass: 2, CodeDetail: 5}

	// Client error codes
	BadRequest               = &ResponseCode{CodeClass: 4, CodeDetail: 0}
	Unauthorized             = &ResponseCode{CodeClass: 4, CodeDetail: 1}
	BadOption                = &ResponseCode{CodeClass: 4, CodeDetail: 2}
	Forbidden                = &ResponseCode{CodeClass: 4, CodeDetail: 3}
	NotFound                 = &ResponseCode{CodeClass: 4, CodeDetail: 4}
	MethodNotAllowed         = &ResponseCode{CodeClass: 4, CodeDetail: 5}
	NotAcceptable            = &ResponseCode{CodeClass: 4, CodeDetail: 6}
	PreconditionFailed       = &ResponseCode{CodeClass: 4, CodeDetail: 12}
	RequestEntityTooLarge    = &ResponseCode{CodeClass: 4, CodeDetail: 13}
	UnsupportedContentFormat = &ResponseCode{CodeClass: 4, CodeDetail: 15}

	// Server error codes
	InternalServerError  = &ResponseCode{CodeClass: 5, CodeDetail: 0}
	NotImplemented       = &ResponseCode{CodeClass: 5, CodeDetail: 1}
	BadGateway           = &ResponseCode{CodeClass: 5, CodeDetail: 2}
	ServiceUnavailable   = &ResponseCode{CodeClass: 5, CodeDetail: 3}
	GatewayTimeout       = &ResponseCode{CodeClass: 5, CodeDetail: 4}
	ProxyingNotSupported = &ResponseCode{CodeClass: 5, CodeDetail: 5}
)

func (r *ResponseCode) ToCodeType() *CodeType {
	return (*CodeType)(r)
}

func (c *CodeType) ToResponseCode() *ResponseCode {
	return (*ResponseCode)(c)
}

func (m *MethodCode) ToCodeType() *CodeType {
	return (*CodeType)(m)
}

func (r *ResponseCode) String() string {
	switch r {
	case Ok:
		return "Ok"
	case Created:
		return "Created"
	case Deleted:
		return "Deleted"
	case Valid:
		return "Valid"
	case Changed:
		return "Changed"
	case Content:
		return "Content"
	case BadRequest:
		return "BadRequest"
	case Unauthorized:
		return "Unauthorized"
	case BadOption:
		return "BadOption"
	case Forbidden:
		return "Forbidden"
	case NotFound:
		return "NotFound"
	case MethodNotAllowed:
		return "MethodNotAllowed"
	case NotAcceptable:
		return "NotAcceptable"
	case PreconditionFailed:
		return "PreconditionFailed"
	case RequestEntityTooLarge:
		return "RequestEntityTooLarge"
	case UnsupportedContentFormat:
		return "UnsupportedContentFormat"
	case InternalServerError:
		return "InternalServerError"
	case NotImplemented:
		return "NotImplemented"
	case BadGateway:
		return "BadGateway"
	case ServiceUnavailable:
		return "ServiceUnavailable"
	case GatewayTimeout:
		return "GatewayTimeout"
	case ProxyingNotSupported:
		return "ProxyingNotSupported"
	default:
		slf4go.GetLogger("codes").Panicf("unsupported response code: %d", r)
		return ""
	}
}
