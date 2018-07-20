package coap

import (
	"fmt"
	"time"
)

// TransmissionParameters which can be adjusted
type TransmissionParameters struct {
	AckTimeout      time.Duration
	AckRandomFactor float64
	MaxRetransmit   int
	NStart          int
	DefaultLeisure  time.Duration
	ProbingRate     int
}

func DefaultTransmissionParameters() TransmissionParameters {
	p := new(TransmissionParameters)

	p.AckTimeout = 2 * time.Second
	p.AckRandomFactor = 1.5
	p.MaxRetransmit = 4
	p.NStart = 1
	p.DefaultLeisure = 5 * time.Second
	p.ProbingRate = 1

	return *p
}

//func ValidateParameters(params *TransmissionParameters) {
//	logger := slf4go.GetLogger("transmission")
//	logger.Debugf("parameters: %v", params)
//}

func (p TransmissionParameters) String() string {
	return fmt.Sprintf("{ AckTimeout: %v, AckRandomFactor: %v, MaxRetransmit: %v, NStart: %v, DefaultLeisure: %v, ProbingRate: %v }",
		p.AckTimeout, p.AckRandomFactor, p.MaxRetransmit, p.NStart, p.DefaultLeisure, p.ProbingRate)
}
