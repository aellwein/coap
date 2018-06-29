package message

import (
	"bytes"
	"fmt"
)

// Token, max 8 bytes.
type TokenType []byte

func (t TokenType) String() string {
	var b bytes.Buffer
	b.WriteString("0x")
	for _, i := range t {
		b.WriteString(fmt.Sprintf("%02X", i))
	}
	return b.String()
}
