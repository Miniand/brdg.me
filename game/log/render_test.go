package log

import (
	"testing"
)

func TestRenderMessage(t *testing.T) {
	m := NewPublicMessage("Test")
	t.Log(RenderMessage(m))
}

func TestRenderMessages(t *testing.T) {
	messages := []Message{
		NewPublicMessage("hello"),
		NewPrivateMessage("Goodbye", []string{"joe"}),
	}
	t.Log(RenderMessages(messages))
}
