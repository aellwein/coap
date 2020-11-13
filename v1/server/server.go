package server

import (
	"githib.com/aellwein/coap/v1"
)

type builder struct {
	c *ctx
}

type ctx struct {
	serveCoap  bool
	serveCoaps bool
	portCoap   uint16
	portCoaps  uint16
}

// Builder creates a new builder for server.
func Builder() *builder {
	return &builder{
		c: &ctx{
			serveCoap:  true,
			serveCoaps: true,
			portCoap:   v1.CoapPort,
			portCoaps:  v1.CoapsPort,
		},
	}
}

func (b *builder) EnableCoap(enable bool) {
	b.c.serveCoap = enable
}

func (b *builder) EnableCoapOnPort(port uint16) {
	b.c.serveCoap = true
	b.c.portCoap = port
}

func (b *builder) EnableCoaps(enable bool) {
	b.c.serveCoaps = enable
}

func (b *builder) EnableCoapsOnPort(port uint16) {
	b.c.serveCoaps = true
	b.c.portCoaps = port
}

func (b *builder) Build() *ctx {
	return b.c
}

func (c *ctx) Listen() error {
	return nil
}
