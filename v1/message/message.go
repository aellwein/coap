package message

func (m *mBuilder) Code(mCode *MessageCode) *mCodeBuilder {
	m.msg.code = mCode
	return &mCodeBuilder{msg: m.msg}
}

func (m *mCodeBuilder) Option(o *Option) *mOptionBuilder {
	m.msg.opts = append(m.msg.opts, o)
	return &mOptionBuilder{msg: m.msg}
}

func (m *mCodeBuilder) Options(o ...*Option) *mOptionBuilder {
	m.msg.opts = append(m.msg.opts, o...)
	return &mOptionBuilder{msg: m.msg}
}

func (m *message) Bytes() []byte {
	return []byte{}
}

func (m *Message) IsEmpty() bool {
	return m.code.IsEmpty()
}

func (m *Message) IsRequest() bool {
	return m.code.IsRequest()
}

func (m *Message) IsResponse() bool {
	return m.code.IsResponse()
}

func (m *Message) Type() MessageType {
	return m.t
}

type Message struct {
	t    MessageType
	code *MessageCode
	opts []*Option
}

type mBuilder struct {
	msg *Message
}

type mCodeBuilder mBuilder
type mOptionBuilder mBuilder
type message mBuilder
