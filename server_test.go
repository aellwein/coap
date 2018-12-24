package coap

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"bou.ke/monkey"
	c "github.com/smartystreets/goconvey/convey"
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

func TestNewInecureCoapServerWithEmptyResourcePath(t *testing.T) {
	c.Convey("Given a new insecure coap server", t, func() {
		c.Convey("When the resource path is empty", func() {
			res := &Resource{Path: ""}
			_, err := NewInsecureCoapServerWithDefaultParameters(res)

			c.Convey("Then an error is returned", func() {
				c.So(err, c.ShouldNotBeNil)
			})
		})
	})
}

func TestNewInecureCoapServerWithResourcePathNotStartingWithSlash(t *testing.T) {
	c.Convey("Given a new insecure coap server", t, func() {
		c.Convey("When the resource path is empty", func() {
			res := &Resource{Path: "rd"}
			_, err := NewInsecureCoapServerWithDefaultParameters(res)

			c.Convey("Then an error is returned", func() {
				c.So(err, c.ShouldNotBeNil)
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

type routingTestCase struct {
	code                 *CodeType
	resource             *Resource
	expectedResponseCode *CodeType
}

var routingTestCases = []routingTestCase{
	{
		code:                 POST,
		resource:             &Resource{Path: "/notExists"},
		expectedResponseCode: NotFound,
	},
	{
		code:                 POST,
		resource:             &Resource{Path: "/rd"},
		expectedResponseCode: MethodNotAllowed,
	},
	{
		code:                 PUT,
		resource:             &Resource{Path: "/rd"},
		expectedResponseCode: MethodNotAllowed,
	},
	{
		code:                 GET,
		resource:             &Resource{Path: "/rd"},
		expectedResponseCode: MethodNotAllowed,
	},
	{
		code:                 DELETE,
		resource:             &Resource{Path: "/rd"},
		expectedResponseCode: MethodNotAllowed,
	},
	{
		code: POST,
		resource: &Resource{Path: "/rd", OnPOST: func(request *Message) (*Message, error) {
			return nil, errors.New("provoked")
		}},
		expectedResponseCode: InternalServerError,
	},
	{
		code: PUT,
		resource: &Resource{Path: "/rd", OnPUT: func(request *Message) (*Message, error) {
			return nil, errors.New("provoked")
		}},
		expectedResponseCode: InternalServerError,
	},
	{
		code: GET,
		resource: &Resource{Path: "/rd", OnGET: func(request *Message) (*Message, error) {
			return nil, errors.New("provoked")
		}},
		expectedResponseCode: InternalServerError,
	},
	{
		code: DELETE,
		resource: &Resource{Path: "/rd", OnDELETE: func(request *Message) (*Message, error) {
			return nil, errors.New("provoked")
		}},
		expectedResponseCode: InternalServerError,
	},
	{
		code: POST,
		resource: &Resource{Path: "/rd", OnPOST: func(request *Message) (*Message, error) {
			return NewContentResponseMessage(request), nil
		}},
		expectedResponseCode: Content,
	},
	{
		code: PUT,
		resource: &Resource{Path: "/rd", OnPUT: func(request *Message) (*Message, error) {
			return NewContentResponseMessage(request), nil
		}},
		expectedResponseCode: Content,
	},
	{
		code: GET,
		resource: &Resource{Path: "/rd", OnGET: func(request *Message) (*Message, error) {
			return NewContentResponseMessage(request), nil
		}},
		expectedResponseCode: Content,
	},
	{
		code: DELETE,
		resource: &Resource{Path: "/rd", OnDELETE: func(request *Message) (*Message, error) {
			return NewContentResponseMessage(request), nil
		}},
		expectedResponseCode: Content,
	},
}

func TestServer_RouteMessage(t *testing.T) {
	for _, tc := range routingTestCases {
		c.Convey(fmt.Sprintf("Given a coap server with resource %v", tc.resource), t, func() {

			server, _ := NewInsecureCoapServerWithDefaultParameters(tc.resource)

			msg := NewConfirmableMessageBuilder().
				Code(tc.code).
				WithRandomMessageId().
				WithRandomToken().
				Option(UriPath, []byte("rd")).
				Build()

			c.Convey(fmt.Sprintf("When a %v message is routed", msg.Code), func() {
				resp := server.routeRequest(msg)

				c.Convey(fmt.Sprintf("Then the response has the code '%v'", tc.expectedResponseCode), func() {
					c.So(*resp.Code, c.ShouldResemble, *tc.expectedResponseCode)
				})
			})
		})
	}
}

func TestServer_RouteMessageWithNoUriPath(t *testing.T) {
	c.Convey("Given a coap server", t, func() {
		server, _ := NewInsecureCoapServerWithDefaultParameters(&Resource{Path: "/rd"})

		c.Convey("When a message without uri-path option is routed", func() {

			msg := NewConfirmableMessageBuilder().
				Code(POST).
				WithRandomMessageId().
				WithRandomToken().
				Build()

			c.Convey("Then a response with BadRequest code is created", func() {
				resp := server.routeRequest(msg)

				c.So(*resp.Code, c.ShouldResemble, *BadRequest)
			})
		})
	})
}
func TestServer_RouteMessageWithInvalidMethodCode(t *testing.T) {
	c.Convey("Given a coap server", t, func() {
		server, _ := NewInsecureCoapServerWithDefaultParameters(&Resource{Path: "/rd"})

		c.Convey("When a message with invalid method code is routed", func() {

			msg := NewConfirmableMessageBuilder().
				Code(EmptyMessage).
				WithRandomMessageId().
				WithRandomToken().
				Option(UriPath, []byte("rd")).
				Build()

			c.Convey("Then a response with BadRequest code is created", func() {
				resp := server.routeRequest(msg)

				c.So(*resp.Code, c.ShouldResemble, *BadRequest)
			})
		})
	})
}
