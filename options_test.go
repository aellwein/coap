package coap

import (
	c "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUriPathOptionToString(t *testing.T) {
	c.Convey("Given Uri-Path option", t, func() {

		uriPath := []OptionValueType{
			[]byte("dead"),
			[]byte("beef"),
			[]byte("cafe"),
			[]byte("babe"),
		}

		c.Convey("When Uri-Path converted to string", func() {
			path := UriPathOptionToString(uriPath)

			c.Convey("Then the result is equal the expected path", func() {
				c.So(path, c.ShouldEqual, "/dead/beef/cafe/babe")
			})
		})
	})
}
