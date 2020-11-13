package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	empty := &MessageCode{0, 0}
	nonEmpty := &MessageCode{0, 1}
	assert.True(t, empty.IsEmpty())
	assert.True(t, EmptyMessage.IsEmpty())
	assert.False(t, nonEmpty.IsEmpty())
}

func TestIsRequest(t *testing.T) {
	nonReq := &MessageCode{1, 10}
	assert.False(t, EmptyMessage.IsRequest())
	assert.False(t, nonReq.IsRequest())
	assert.True(t, GET.IsRequest())
}

func TestIsResponse(t *testing.T) {
	resp := &MessageCode{2, 0}
	assert.False(t, EmptyMessage.IsResponse())
	assert.True(t, resp.IsResponse())
	assert.True(t, BadRequest.IsResponse())
}
