package transmission

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCopyFrom(t *testing.T) {
	Convey("Given a set of parameters", t, func() {
		p := Parameters{
			AckTimeout:      2,
			AckRandomFactor: 3,
			MaxRetransmit:   10,
			NStart:          6,
			ProbingRate:     1,
			DefaultLeisure:  4,
		}
		Convey("When copied", func() {
			p2 := CopyFrom(p)
			Convey("Then parameters should be the same", func() {
				So(*p2, ShouldResemble, p)
			})
		})
	})
}

func TestParameters_String(t *testing.T) {
	t.Log(NewDefaultParameters())
}
