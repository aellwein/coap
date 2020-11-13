package message

type reqBuilder struct {
	msg *Message
}
type typedRequestBuilder reqBuilder

// RequestBuilder creates a new request builder
func RequestBuilder() *reqBuilder {
	return &reqBuilder{}
}

func (*reqBuilder) newRequestOfType(mType MessageType, code *MessageCode) *typedRequestBuilder {
	return &typedRequestBuilder{
		msg: &Message{
			t:    mType,
			code: code,
		},
	}
}

func (r *reqBuilder) NewConfirmableGET() *typedRequestBuilder {
	return r.newRequestOfType(CON, GET)
}

func (r *reqBuilder) NewNonConfirmableGET() *typedRequestBuilder {
	return r.newRequestOfType(NON, GET)
}

func (r *reqBuilder) NewConfirmablePUT() *typedRequestBuilder {
	return r.newRequestOfType(CON, PUT)
}

func (r *reqBuilder) NewNonConfirmablePUT() *typedRequestBuilder {
	return r.newRequestOfType(NON, PUT)
}

func (r *reqBuilder) NewConfirmablePOST() *typedRequestBuilder {
	return r.newRequestOfType(CON, POST)
}

func (r *reqBuilder) NewNonConfirmablePOST() *typedRequestBuilder {
	return r.newRequestOfType(NON, POST)
}

func (r *reqBuilder) NewConfirmableDELETE() *typedRequestBuilder {
	return r.newRequestOfType(CON, DELETE)
}

func (r *reqBuilder) NewNonConfirmableDELETE() *typedRequestBuilder {
	return r.newRequestOfType(NON, DELETE)
}

func (r *reqBuilder) NewEmptyConfirmable() *typedRequestBuilder {
	return r.newRequestOfType(CON, EmptyMessage)
}

func (r *reqBuilder) NewEmptyNonConfirmable() *typedRequestBuilder {
	return r.newRequestOfType(NON, EmptyMessage)
}
