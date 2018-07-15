package coap

import "fmt"

type ContentType uint16

const (
	ContentTypeTextPlain              ContentType = iota
	ContentTypeApplicationLinkFormat              = 40
	ContentTypeApplicationXml                     = 41
	ContentTypeApplicationOctetStream             = 42
	ContentTypeApplicationExi                     = 47
	ContentTypeApplicationJson                    = 50
)

var AllContentTypes = []ContentType{
	ContentTypeTextPlain,
	ContentTypeApplicationLinkFormat,
	ContentTypeApplicationXml,
	ContentTypeApplicationOctetStream,
	ContentTypeApplicationExi,
	ContentTypeApplicationJson,
}

func (c ContentType) String() string {
	switch c {
	case ContentTypeTextPlain:
		return fmt.Sprintf("%d (%s)", c, "text/plain;charset=utf-8")
	case ContentTypeApplicationLinkFormat:
		return fmt.Sprintf("%d (%s)", c, "application/link-format")
	case ContentTypeApplicationXml:
		return fmt.Sprintf("%d (%s)", c, "application/xml")
	case ContentTypeApplicationOctetStream:
		return fmt.Sprintf("%d (%s)", c, "application/octet-stream")
	case ContentTypeApplicationExi:
		return fmt.Sprintf("%d (%s)", c, "application/exi")
	case ContentTypeApplicationJson:
		return fmt.Sprintf("%d (%s)", c, "application/json")
	default:
		return fmt.Sprintf("%d (%s)", c, "unknown")
	}
}
