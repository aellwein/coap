package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
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

func NewToken() *TokenType {
	t := TokenType(make([]byte, 8))
	binary.BigEndian.PutUint64(t, rand.Uint64())
	return &t
}

func (t *TokenType) Copy() *TokenType {
	to := make([]byte, 8)
	copy(to, *t)
	return t
}
