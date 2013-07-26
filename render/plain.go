package render

import (
	"bytes"
	"github.com/beefsack/brdg.me/game"
	"html/template"
)

type PlainMarkupper struct {
	Markupper
}

func (t *PlainMarkupper) StartColour(colour string) interface{} {
	return ""
}
func (t *PlainMarkupper) EndColour() interface{} {
	return ""
}
func (t *PlainMarkupper) StartBold() interface{} {
	return ""
}
func (t *PlainMarkupper) EndBold() interface{} {
	return ""
}

func RenderPlain(tmpl string, g game.Playable) (string, error) {
	t := template.Must(template.New("tmpl").
		Funcs(AttachTemplateFuncs(template.FuncMap{}, &PlainMarkupper{})).Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, g)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
