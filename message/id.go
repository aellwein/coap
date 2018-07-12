package message

import (
	"fmt"
	"math/rand"
)

// Message ID, 16 bits.
type MessageIdType uint16

func (m MessageIdType) String() string {
	return fmt.Sprintf("0x%04X", uint16(m))
}

func NewMessageId() MessageIdType {
	return MessageIdType(rand.Uint32())
}
