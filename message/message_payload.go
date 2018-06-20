package message

import (
	"bytes"
	"fmt"
)

type PayloadType []byte

func (p PayloadType) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	for n, i := range p {
		if n == 0 {
			b.WriteString(fmt.Sprintf("0x%02X", i))
		} else {
			b.WriteString(fmt.Sprintf(", 0x%02X", i))
		}
	}
	b.WriteString("]")
	return b.String()
}
