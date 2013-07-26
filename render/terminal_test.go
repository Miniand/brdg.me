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
	if output != "\x1b[34m\x1b[34;1mhello\x1b[34m\x1b[0m" {
		t.Error("Output was", output)
		return
	}
}
