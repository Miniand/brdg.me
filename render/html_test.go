package render

import (
	"testing"
)

func TestHtmlRender(t *testing.T) {
	output, err := RenderHtml(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`)
	if err != nil {
		t.Error(err)
		return
	}
	if output != `<span style="color:rgb(25,118,210);"><strong>hello</strong></span>` {
		t.Error("Output was", output)
		return
	}
}
