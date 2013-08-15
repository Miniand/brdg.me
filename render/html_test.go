package render

import (
	"testing"
)

func TestHtmlRender(t *testing.T) {
	output, err := RenderHtml(`{{l}}{{c "blue"}}{{b}}hello{{_b}}{{_c}}{{_l}}`)
	if err != nil {
		t.Error(err)
		return
	}
	if output != `<span style="font-size:1.6em;"><span style="color:rgb(0,0,187);"><b>hello</b></span></span>` {
		t.Error("Output was", output)
		return
	}
}
