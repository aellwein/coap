package message

import "github.com/aellwein/coap/logging"

// Basic message type, fits in 2 bits.
type MessageType int8

// Message types used in CoAP.
const (
	CON MessageType = iota
	NON
	ACK
	RST
)

// message type to string
func (mt MessageType) String() string {
	switch mt {
	case CON:
		return "CON"
	case NON:
		return "NON"
	case ACK:
		return "ACK"
	case RST:
		return "RST"
	default:
		logging.Sugar.Panicf("unknown message type: %d", mt)
		return ""
	}
}