package coap

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aellwein/slf4go"
	_ "github.com/aellwein/slf4go-logrus-adaptor"
	"net"
	"strings"
)

type CoapPort uint16

const (
	InsecurePort CoapPort = 5683
	SecurePort   CoapPort = 5684
)
const MaxPacketSize = 2048

type resourceMap map[string]*Resource

type Server struct {
	addr       *net.UDPAddr
	conn       *net.UDPConn
	parameters TransmissionParameters
	resources  resourceMap
}

var logger slf4go.Logger

// Get string representation of the server
func (server Server) String() string {
	return fmt.Sprintf("Server{ addr=%v, parameters=%v, conn=%v, resources=%v}",
		server.addr, server.parameters, server.conn, server.resources)
}

func (r resourceMap) String() string {
	var s bytes.Buffer
	m := make([]string, 0)
	for k, v := range r {
		s.WriteString(fmt.Sprintf("{ '%v': [", k))
		if (*v).OnPOST != nil {
			m = append(m, "POST")
		}
		if (*v).OnGET != nil {
			m = append(m, "GET")
		}
		if (*v).OnPUT != nil {
			m = append(m, "PUT")
		}
		if (*v).OnDELETE != nil {
			m = append(m, "DELETE")
		}
		s.WriteString(strings.Join(m, ","))
		s.WriteString("]}")
	}
	return s.String()
}

func newServer(port CoapPort, parameters TransmissionParameters, resources ...*Resource) (*Server, error) {
	var err error
	server := &Server{}

	logger = slf4go.GetLogger("server")
	server.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	//transmission.ValidateParameters(parameters)

	server.parameters = parameters
	server.resources = make(map[string]*Resource)

	for _, r := range resources {
		if r.Path == "" || r.Path[0:1] != "/" {
			return nil, errors.New("path may not be empty and must start with slash")
		}
		server.resources[r.Path] = r
	}

	return server, nil
}

// Creates a default CoAP Server on secure port using default transmission parameters.
func NewSecureCoapServerWithDefaultParameters(resources ...*Resource) (*Server, error) {
	return newServer(SecurePort, DefaultTransmissionParameters(), resources...)
}

// Creates a default CoAP server on insecure port using default transmission parameters.
func NewInsecureCoapServerWithDefaultParameters(resources ...*Resource) (*Server, error) {
	return newServer(InsecurePort, DefaultTransmissionParameters(), resources...)
}

// Creates a new CoAP server on secure port using given transmission parameters.
func NewSecureCoapServer(parameters TransmissionParameters, resources ...*Resource) (*Server, error) {
	params := parameters
	return newServer(SecurePort, params, resources...)
}

// Creates a default CoAP server on insecure port using given transmission parameters.
func NewInsecureCoapServer(parameters TransmissionParameters, resources ...*Resource) (*Server, error) {
	params := parameters
	return newServer(InsecurePort, params, resources...)
}

// Listen on default (insecure) port
func (server *Server) Listen() error {
	return server.ListenOn(InsecurePort)
}

func (s *Server) handlePacket(packet []byte, n int, peer *net.UDPAddr) {
	logger.Debugf("received packet from %s: \n%s", peer, hex.Dump(packet[0:n]))
	msg, err := NewMessageFromBytesAndPeer(packet[0:n], peer)
	if err != nil {
		logger.Debugf("error decoding message: %v", err)
		// message could not be decoded, ignore
		return
	}
	logger.Debugf("message received: %v", msg)
	logger.Debug("Go representation of the packet: ", DumpInGoFormat(packet[0:n]))

	if msg.Type == NonConfirmable || msg.Type == Confirmable {

		if res := msg.Validate(); res != Ok {
			b := responseWithCode(msg, res).ToBytes()
			s.conn.WriteToUDP(b, peer)
			return
		}

		// route request and get response
		resp := s.routeRequest(msg)

		logger.Debugf("will send message %v", resp)
		// write response
		respBuf := resp.ToBytes()
		s.conn.WriteToUDP(respBuf, peer)
	}

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
	logger.Infof("Server is listening on %v", server.conn.LocalAddr())

	for {
		n, peer, err := server.conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Debug(err)
			// try to read again if read failed
			continue
		}

		server.handlePacket(buffer, n, peer)
	}
}

func (server *Server) routeRequest(msg *Message) *Message {
	if pathOption, ok := (*msg.Options)[UriPath]; ok {
		p := UriPathOptionToString(pathOption)
		if handler, ok := server.resources[p]; ok {

			switch *msg.Code {

			case *GET:
				if handler.OnGET != nil {
					if resp, err := handler.OnGET(msg); err != nil {
						return NewInternalServerErrorResponseMessage(msg)
					} else {
						return resp
					}
				}

			case *POST:
				if handler.OnPOST != nil {
					if resp, err := handler.OnPOST(msg); err != nil {
						return NewInternalServerErrorResponseMessage(msg)
					} else {
						return resp
					}
				} else {
					return NewMethodNotAllowedResponseMessage(msg)
				}

			case *PUT:
				if handler.OnPUT != nil {
					if resp, err := handler.OnPUT(msg); err != nil {
						return NewInternalServerErrorResponseMessage(msg)
					} else {
						return resp
					}
				} else {
					return NewMethodNotAllowedResponseMessage(msg)
				}

			case *DELETE:
				if handler.OnDELETE != nil {
					if resp, err := handler.OnDELETE(msg); err != nil {
						return NewInternalServerErrorResponseMessage(msg)
					} else {
						return resp
					}
				} else {
					return NewMethodNotAllowedResponseMessage(msg)
				}

			default:
				return NewBadRequestResponseMessage(msg)
			}

		} else {
			// no handler found
			return NewNotFoundResponseMessage(msg)
		}
	} else {
		// no path in message, bad request
		return NewBadRequestResponseMessage(msg)
	}
	// should never happen
	return NewInternalServerErrorResponseMessage(msg)
}

func (s *Server) AddResource(resource *Resource) {
	s.resources[resource.Path] = resource
}

func (s *Server) RemoveResource(resource *Resource) {
	if _, exists := s.resources[resource.Path]; exists {
		delete(s.resources, resource.Path)
	}
}

func (s *Server) RemoveResourceByPath(path string) {
	if _, exists := s.resources[path]; exists {
		delete(s.resources, path)
	}
}
