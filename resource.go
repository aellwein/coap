package coap

// ResourceHandlerFunc is a shortcut type for resource handling function.
type ResourceHandlerFunc func(request *Request) (Response, error)

type Resource struct {
	Path     string
	OnGET    ResourceHandlerFunc
	OnPUT    ResourceHandlerFunc
	OnPOST   ResourceHandlerFunc
	OnDELETE ResourceHandlerFunc
}

var defaultResourceHandlerFunc ResourceHandlerFunc = func(request *Request) (Response, error) {
	return NewMethodNotAllowedResponse(request), nil
}

func NewResource(path string) *Resource {
	return &Resource{
		Path:     path,
		OnGET:    defaultResourceHandlerFunc,
		OnPUT:    defaultResourceHandlerFunc,
		OnPOST:   defaultResourceHandlerFunc,
		OnDELETE: defaultResourceHandlerFunc,
	}
}
