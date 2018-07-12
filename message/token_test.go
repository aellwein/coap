package message

import (
	c "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTokenType_Copy(t *testing.T) {
	c.Convey("Given a pre-defined token", t, func() {
		tkn := TokenType([]byte{0x13, 0x37, 0xCA, 0xFE})
		c.Convey("When token is copied", func() {
			tkn2 := *tkn.Copy()
			c.Convey("Then the resulting token matches exactly the given one", func() {
				c.So(tkn2, c.ShouldResemble, tkn)
			})
		})
	})
}
