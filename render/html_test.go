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
	if output != `<span style="color:rgb(0,0,187);"><b>hello</b></span>` {
		t.Error("Output was", output)
		return
	}
}
