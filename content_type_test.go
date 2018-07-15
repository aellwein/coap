package coap

import "testing"

func TestContentType_String(t *testing.T) {
	for _, c := range AllContentTypes {
		s := c.String()
		if s == "" {
			t.Error("content type may not be empty")
		}
	}
	if ContentType(65521).String() == "" {
		t.Error("unknown content type may not be empty")
	}
}
