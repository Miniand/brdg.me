package render

import (
	"testing"
)

func TestTerminalRender(t *testing.T) {
	output, err := RenderTerminal(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`)
	if err != nil {
		t.Error(err)
		return
	}
	if output != "\x1b[0;34;49m\x1b[0;34;49;1mhello\x1b[0;34;49m\x1b[0;39;49m" {
		t.Error("Output was", output)
		return
	}
}
