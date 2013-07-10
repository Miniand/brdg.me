package render

import (
	"github.com/beefsack/boredga.me/game/tic_tac_toe"
	"testing"
)

func TestHtmlRender(t *testing.T) {
	g := &tic_tac_toe.Game{}
	output, err := RenderHtml(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`, g)
	if err != nil {
		t.Error(err)
		return
	}
	if output != `<span style="color:rgb(0,0,187);"><b>hello</b></span>` {
		t.Error("Output was", output)
		return
	}
}
