package render

import (
	"testing"
)

func TestPlainRender(t *testing.T) {
	output := RenderPlain(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`)
	if output != `hello` {
		t.Error("Output was", output)
		return
	}
}
