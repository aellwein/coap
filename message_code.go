package coap

import (
	"fmt"
)

/* MESSAGE CODES */

// Empty message code
var EmptyMessage = &CodeType{CodeClass: 0, CodeDetail: 0}

// Method Codes
var (
	GET    = &CodeType{CodeClass: 0, CodeDetail: 1}
	POST   = &CodeType{CodeClass: 0, CodeDetail: 2}
	PUT    = &CodeType{CodeClass: 0, CodeDetail: 3}
	DELETE = &CodeType{CodeClass: 0, CodeDetail: 4}
)

var (
	// success codes
	Ok      = &CodeType{CodeClass: 2, CodeDetail: 0}
	Created = &CodeType{CodeClass: 2, CodeDetail: 1}
	Deleted = &CodeType{CodeClass: 2, CodeDetail: 2}
	Valid   = &CodeType{CodeClass: 2, CodeDetail: 3}
	Changed = &CodeType{CodeClass: 2, CodeDetail: 4}
	Content = &CodeType{CodeClass: 2, CodeDetail: 5}

	// Client error codes
	BadRequest               = &CodeType{CodeClass: 4, CodeDetail: 0}
	Unauthorized             = &CodeType{CodeClass: 4, CodeDetail: 1}
	BadOption                = &CodeType{CodeClass: 4, CodeDetail: 2}
	Forbidden                = &CodeType{CodeClass: 4, CodeDetail: 3}
	NotFound                 = &CodeType{CodeClass: 4, CodeDetail: 4}
	MethodNotAllowed         = &CodeType{CodeClass: 4, CodeDetail: 5}
	NotAcceptable            = &CodeType{CodeClass: 4, CodeDetail: 6}
	PreconditionFailed       = &CodeType{CodeClass: 4, CodeDetail: 12}
	RequestEntityTooLarge    = &CodeType{CodeClass: 4, CodeDetail: 13}
	UnsupportedContentFormat = &CodeType{CodeClass: 4, CodeDetail: 15}

	// Server error codes
	InternalServerError  = &CodeType{CodeClass: 5, CodeDetail: 0}
	NotImplemented       = &CodeType{CodeClass: 5, CodeDetail: 1}
	BadGateway           = &CodeType{CodeClass: 5, CodeDetail: 2}
	ServiceUnavailable   = &CodeType{CodeClass: 5, CodeDetail: 3}
	GatewayTimeout       = &CodeType{CodeClass: 5, CodeDetail: 4}
	ProxyingNotSupported = &CodeType{CodeClass: 5, CodeDetail: 5}
)

var AllMessageCodes = []*CodeType{
	EmptyMessage,
	GET,
	POST,
	PUT,
	DELETE,
	Ok,
	Created,
	Deleted,
	Valid,
	Changed,
	Content,
	BadRequest,
	Unauthorized,
	BadOption,
	Forbidden,
	NotFound,
	MethodNotAllowed,
	NotAcceptable,
	PreconditionFailed,
	RequestEntityTooLarge,
	UnsupportedContentFormat,
	InternalServerError,
	NotImplemented,
	BadGateway,
	ServiceUnavailable,
	GatewayTimeout,
	ProxyingNotSupported,
}

func (c *CodeType) String() string {
	switch *c {
	case *EmptyMessage:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Empty")
	case *GET:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "GET")
	case *POST:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "POST")
	case *PUT:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "PUT")
	case *DELETE:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "DELETE")
	case *Ok:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Ok")
	case *Created:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Created")
	case *Deleted:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Deleted")
	case *Valid:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Valid")
	case *Changed:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Changed")
	case *Content:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Content")
	case *BadRequest:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "BadRequest")
	case *Unauthorized:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Unauthorized")
	case *BadOption:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "BadOption")
	case *Forbidden:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Forbidden")
	case *NotFound:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "NotFound")
	case *MethodNotAllowed:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "MethodNotAllowed")
	case *NotAcceptable:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "NotAcceptable")
	case *PreconditionFailed:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "PreconditionFailed")
	case *RequestEntityTooLarge:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "RequestEntityTooLarge")
	case *UnsupportedContentFormat:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "UnsupportedContentFormat")
	case *InternalServerError:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "InternalServerError")
	case *NotImplemented:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "NotImplemented")
	case *BadGateway:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "BadGateway")
	case *ServiceUnavailable:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "ServiceUnavailable")
	case *GatewayTimeout:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "GatewayTimeout")
	case *ProxyingNotSupported:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "ProxyingNotSupported")
	default:
		return fmt.Sprintf("%d.%02d (%s)", c.CodeClass, c.CodeDetail, "Unknown Code")
	}
}
