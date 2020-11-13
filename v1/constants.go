package v1

import "time"

const (
	// CoapPort defines default port to listen for CoAP requests
	CoapPort uint16 = 5683
	// CoapsPort defines default port to listen for CoAPS requests
	CoapsPort uint16 = 5684

	// AckTimeout
	AckTimeout = 2 * time.Second
	// AckRandomFactor
	AckRandomFactor = 1.5
	// MaxRetransmit
	MaxRetransmit = 4
	// NStart
	NStart = 1
	// DefaultLeisure
	DefaultLeisure = 5 * time.Second
	// ProbingRate
	ProbingRate = 1

	MessageVersion = 0x01
)
