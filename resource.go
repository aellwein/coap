package coap

import "fmt"

// ResourceHandlerFunc is a shortcut type for resource handling function.
type ResourceHandlerFunc func(request *Message) (*Message, error)

type Resource struct {
	Path     string
	OnGET    ResourceHandlerFunc
	OnPUT    ResourceHandlerFunc
	OnPOST   ResourceHandlerFunc
	OnDELETE ResourceHandlerFunc
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource{ Path: '%v', OnGET: %v, OnPUT: %v, OnPOST: %v, OnDELETE: %v }",
		r.Path, r.OnGET, r.OnPUT, r.OnPOST, r.OnDELETE)
}
