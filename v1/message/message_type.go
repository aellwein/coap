package message

// MessageType indicates the type of message
type MessageType uint8

const (
	// CON represents a Confirmable message type
	CON MessageType = iota
	// NON represents Non-Confirmable message type
	NON
	// ACK represents Acknowledgement message type
	ACK
	// RST represents Reset message type
	RST
)

func (m MessageType) String() string {
	switch m {
	case CON:
		return "CON"
	case NON:
		return "NON"
	case ACK:
		return "ACK"
	case RST:
		return "RST"
	default:
		return "unknown"
	}
}
