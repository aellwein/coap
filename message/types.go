package message

import "github.com/aellwein/coap/logging"

// Basic message type, fits in 2 bits.
type MessageType int8

// Message types used in CoAP.
const (
	Confirmable MessageType = iota
	NonConfirmable
	Acknowledgement
	Reset
)

// message type to string
func (mt MessageType) String() string {
	switch mt {
	case Confirmable:
		return "CON"
	case NonConfirmable:
		return "NON"
	case Acknowledgement:
		return "ACK"
	case Reset:
		return "RST"
	default:
		logging.Sugar.Panicf("unknown message type: %d", mt)
		return ""
	}
}
