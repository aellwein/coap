package coap

import (
	"errors"
	"github.com/bouk/monkey"
	c "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestNewInsecureCoapServer(t *testing.T) {
	c.Convey("Given a new insecure coap server", t, func() {
		c.Convey("When the server is created", func() {
			params := DefaultTransmissionParameters()
			res := &Resource{
				Path: "/",
			}
			server, err := NewInsecureCoapServer(params, res)
			c.Convey("Then the parameters should be the same", func() {
				c.So(err, c.ShouldBeNil)
				c.So(server.String(), c.ShouldNotBeNil)
				c.So(server.parameters, c.ShouldResemble, params)
				c.So(server.addr.Port, c.ShouldEqual, InsecurePort)
			})

			c.Convey("And resources are added to server", func() {
				c.So(server.resources["/"], c.ShouldEqual, res)
			})
		})

	})
}

func TestNewSecureCoapServer(t *testing.T) {
	c.Convey("Given a new secure coap server", t, func() {
		c.Convey("When the server is created", func() {
			params := DefaultTransmissionParameters()
			res := &Resource{
				Path: "/",
			}
			server, err := NewSecureCoapServer(params, res)
			c.Convey("Then the parameters should be the same", func() {
				c.So(err, c.ShouldBeNil)
				c.So(server.String(), c.ShouldNotBeNil)
				c.So(server.parameters, c.ShouldResemble, params)
				c.So(server.addr.Port, c.ShouldEqual, SecurePort)
			})

			c.Convey("And resources are added to server", func() {
				c.So(server.resources["/"], c.ShouldEqual, res)
			})
		})

	})
}

func TestNewSecureCoapServerWithDefaultParameters(t *testing.T) {
	c.Convey("Given a new secure coap server with default parameters", t, func() {
		c.Convey("When the server is created", func() {

			res := &Resource{Path: "/rd"}
			server, err := NewSecureCoapServerWithDefaultParameters(res)

			c.Convey("Then the parameters are equal the default parameters", func() {
				c.So(err, c.ShouldBeNil)
				c.So(server.parameters, c.ShouldResemble, DefaultTransmissionParameters())
				c.So(server.addr.Port, c.ShouldEqual, SecurePort)
			})

			c.Convey("And resources are added to server", func() {
				c.So(server.resources["/rd"], c.ShouldEqual, res)
			})

		})
	})
}

func TestNewInecureCoapServerWithDefaultParameters(t *testing.T) {
	c.Convey("Given a new insecure coap server with default parameters", t, func() {
		c.Convey("When the server is created", func() {

			res := &Resource{Path: "/rd"}
			server, err := NewInsecureCoapServerWithDefaultParameters(res)

			c.Convey("Then the parameters are equal the default parameters", func() {
				c.So(err, c.ShouldBeNil)
				c.So(server.parameters, c.ShouldResemble, DefaultTransmissionParameters())
				c.So(server.addr.Port, c.ShouldEqual, InsecurePort)
			})

			c.Convey("And resources are added to server", func() {
				c.So(server.resources["/rd"], c.ShouldEqual, res)
			})

		})
	})
}

func TestServer_Listen(t *testing.T) {
	c.Convey("Given a new coap server", t, func() {
		server, _ := NewInsecureCoapServerWithDefaultParameters(&Resource{Path: "/rd"})

		mockListenUDP := func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
			panic("mockListenUDP")
		}
		patch := monkey.Patch(net.ListenUDP, mockListenUDP)
		defer patch.Unpatch()

		c.Convey("When Listen is called", func() {

			c.Convey("Then the net.ListenUDP interface is called", func() {
				c.So(func() { server.Listen() }, c.ShouldPanic)
			})
		})
	})
}

func TestServer_ListenWithError(t *testing.T) {
	c.Convey("Given a new coap server", t, func() {
		server, _ := NewInsecureCoapServerWithDefaultParameters(&Resource{Path: "/rd"})

		err := errors.New("Provoked error")

		mockListenUDP := func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
			return nil, err
		}
		patch := monkey.Patch(net.ListenUDP, mockListenUDP)
		defer patch.Unpatch()

		c.Convey("When Listen is called", func() {

			err2 := server.Listen()

			c.Convey("Then the error from underlying interface is passed to user", func() {
				c.So(err2, c.ShouldEqual, err)
			})
		})
	})
}
