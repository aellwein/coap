package coap

import (
	"bytes"
	"fmt"
	"strings"
)

// HexContent makes a hex string representation of a byte array.
func HexContent(p []byte) string {
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

// DumpInGoFormat dumps a byte array in Golang representation.
func DumpInGoFormat(b []byte) string {
	var builder strings.Builder
	builder.WriteString("[]byte{\n")
	for n, i := range b {
		builder.WriteString(fmt.Sprintf(" 0x%02X,", i))
		if (n+1)%16 == 0 {
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n}")
	return builder.String()
}
