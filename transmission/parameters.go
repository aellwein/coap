package transmission

import (
	"fmt"
	"github.com/aellwein/slf4go"
	"time"
)

// Parameters which can be adjusted
type Parameters struct {
	AckTimeout      time.Duration
	AckRandomFactor float64
	MaxRetransmit   int
	NStart          int
	DefaultLeisure  time.Duration
	ProbingRate     int
}

func NewDefaultParameters() *Parameters {
	p := new(Parameters)

	p.AckTimeout = 2 * time.Second
	p.AckRandomFactor = 1.5
	p.MaxRetransmit = 4
	p.NStart = 1
	p.DefaultLeisure = 5 * time.Second
	p.ProbingRate = 1

	return p
}

func CopyFrom(src Parameters) *Parameters {
	p2 := new(Parameters)

	p2.AckTimeout = src.AckTimeout
	p2.AckRandomFactor = src.AckRandomFactor
	p2.MaxRetransmit = src.MaxRetransmit
	p2.NStart = src.NStart
	p2.DefaultLeisure = src.DefaultLeisure
	p2.ProbingRate = src.ProbingRate

	return p2
}

func ValidateParameters(params *Parameters) {
	logger := slf4go.GetLogger("transmission")
	logger.Debugf("parameters: %v", params)
}

func (p *Parameters) String() string {
	return fmt.Sprintf("Parameters{ AckTimeout: %v, AckRandomFactor: %v, MaxRetransmit: %v, NStart: %v, DefaultLeisure: %v, ProbingRate: %v }",
		p.AckTimeout, p.AckRandomFactor, p.MaxRetransmit, p.NStart, p.DefaultLeisure, p.ProbingRate)
}
