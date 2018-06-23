package coap

import (
	"encoding/hex"
	"fmt"
	"github.com/aellwein/coap/logging"
	"github.com/aellwein/coap/message"
	"net"
	"strings"
)

type CoapPort uint16

const (
	InsecurePort CoapPort = 5683
	SecurePort   CoapPort = 5684
)
const MaxPacketSize = 2048

type Server struct {
	addr     *net.UDPAddr
	conn     *net.UDPConn
	handlers map[string]RequestHandler
}

// Get string representation of the server
func (server Server) String() string {
	return fmt.Sprintf("Server{ addr=%v, conn=%v, handlers=%v}",
		server.addr, server.conn, server.handlers)
}

func newServer(port CoapPort, handlers ...RequestHandler) (*Server, error) {
	var err error
	server := &Server{}

	server.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server.handlers = make(map[string]RequestHandler)

	for _, h := range handlers {
		if h.Path != "" {
			server.handlers[h.Path] = h
		}
	}

	return server, nil
}

// default CoAP Server on secure port.
func NewSecureCoapServer(handlers ...RequestHandler) (*Server, error) {
	return newServer(SecurePort, handlers...)
}

// default CoAP server on insecure port.
func NewInsecureCoapServer(handlers ...RequestHandler) (*Server, error) {
	return newServer(InsecurePort, handlers...)
}

// Listen on default (insecure) port
func (server *Server) Listen() error {
	return server.ListenOn(InsecurePort)
}

// Listen on specific port
func (server *Server) ListenOn(port CoapPort) error {
	var err error
	server.conn, err = net.ListenUDP("udp", server.addr)

	if err != nil {
		return err
	}

	defer server.conn.Close()

	buffer := make([]byte, MaxPacketSize)

	for {
		n, peer, err := server.conn.ReadFromUDP(buffer)
		if err != nil {
			logging.Sugar.Debug(err)
		}
		logging.Sugar.Debugf("received packet from %s: \n%s", peer, hex.Dump(buffer[0:n]))
		msg, err := message.DecodeMessage(buffer[0:n], peer)
		if err != nil {
			logging.Sugar.Debugf("error decoding message: %v", err)
		}
		logging.Sugar.Debugf("message received: %v", msg)
		logging.Sugar.Debugf(dumpEncoded(buffer[0:n]))
	}
}

func dumpEncoded(b []byte) string {
	var builder strings.Builder
	builder.WriteString("[]byte{\n")
	for n, i := range b {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("0x%02X,", i))
		if (n+1)%16 == 0 {
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n}")
	return builder.String()
}
