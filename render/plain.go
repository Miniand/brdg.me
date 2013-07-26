package render

import (
	"bytes"
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

func RenderPlain(tmpl string) (string, error) {
	t := template.Must(template.New("tmpl").
		Funcs(AttachTemplateFuncs(template.FuncMap{}, &PlainMarkupper{})).Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, Context{})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
