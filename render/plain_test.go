package render

import (
	"testing"
)

func TestPlainRender(t *testing.T) {
	output, err := RenderPlain(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`)
	if err != nil {
		t.Error(err)
		return
	}
	if output != `hello` {
		t.Error("Output was", output)
		return
	}
}
