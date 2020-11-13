package server

import (
	"testing"

	"githib.com/aellwein/coap/v1"
	"github.com/stretchr/testify/assert"
)

func TestBuilderDefault(t *testing.T) {
	c := Builder().Build()
	assert.Equal(t, v1.CoapPort, c.portCoap)
	assert.Equal(t, v1.CoapsPort, c.portCoaps)
	assert.True(t, c.serveCoap)
	assert.True(t, c.serveCoaps)
}
