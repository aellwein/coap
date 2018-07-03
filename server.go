package coap

import (
	"encoding/hex"
	"fmt"
	"github.com/aellwein/coap/logging"
	"github.com/aellwein/coap/message"
	"github.com/aellwein/coap/transmission"
	"github.com/aellwein/slf4go"
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
	addr       *net.UDPAddr
	conn       *net.UDPConn
	parameters *transmission.Parameters
	handlers   map[string]RequestHandler
}

var logger slf4go.Logger

// Get string representation of the server
func (server Server) String() string {
	return fmt.Sprintf("Server{ addr=%v, parameters=%v, conn=%v, handlers=%v}",
		server.addr, server.parameters, server.conn, server.handlers)
}

func newServer(port CoapPort, parameters *transmission.Parameters, handlers ...RequestHandler) (*Server, error) {
	var err error
	server := &Server{}

	logger = logging.LoggerFactory.GetLogger("server")
	server.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	transmission.ValidateParameters(parameters)

	server.parameters = parameters
	server.handlers = make(map[string]RequestHandler)

	for _, h := range handlers {
		if h.Path != "" {
			server.handlers[h.Path] = h
		}
	}

	return server, nil
}

// Creates a default CoAP Server on secure port using default transmission parameters.
func NewSecureCoapServerWithDefaultParameters(handlers ...RequestHandler) (*Server, error) {
	return newServer(SecurePort, transmission.NewDefaultParameters(), handlers...)
}

// Creates a default CoAP server on insecure port using default transmission parameters.
func NewInsecureCoapServerWithDefaultParameters(handlers ...RequestHandler) (*Server, error) {
	return newServer(InsecurePort, transmission.NewDefaultParameters(), handlers...)
}

// Creates a new CoAP server on secure port using given transmission parameters.
func NewSecureCoapServer(parameters *transmission.Parameters, handlers ...RequestHandler) (*Server, error) {
	params := transmission.CopyFrom(*parameters)
	return newServer(SecurePort, params, handlers...)
}

// Creates a default CoAP server on insecure port using given transmission parameters.
func NewInsecureCoapServer(parameters *transmission.Parameters, handlers ...RequestHandler) (*Server, error) {
	params := transmission.CopyFrom(*parameters)
	return newServer(InsecurePort, params, handlers...)
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
			logger.Debug(err)
		}
		logger.DebugF("received packet from %s: \n%s", peer, hex.Dump(buffer[0:n]))
		msg, err := message.Decode(buffer[0:n], peer)
		if err != nil {
			logger.DebugF("error decoding message: %v", err)
		}
		logger.DebugF("message received: %v", msg)
		logger.Debug("Go representation of the packet: ", dumpEncoded(buffer[0:n]))
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
