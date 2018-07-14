package message

import (
	"math/rand"
	"time"
)

func init() {
	// Needed for Token/Message ID generation
	rand.Seed(time.Now().UnixNano())
}
