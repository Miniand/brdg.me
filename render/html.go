package render

import (
	"bytes"
	"html/template"
)

type HtmlMarkupper struct {
	Markupper
}

func (t *HtmlMarkupper) StartColour(colour string) interface{} {
	return template.HTML(`<span style="color:` + HtmlColours()[colour] + `;">`)
}
func (t *HtmlMarkupper) EndColour() interface{} {
	return template.HTML("</span>")
}
func (t *HtmlMarkupper) StartBold() interface{} {
	return template.HTML("<b>")
}
func (t *HtmlMarkupper) EndBold() interface{} {
	return template.HTML("</b>")
}
func (t *HtmlMarkupper) StartLarge() interface{} {
	return template.HTML(`<span style="font-size:1.6em;">`)
}
func (t *HtmlMarkupper) EndLarge() interface{} {
	return template.HTML("</span>")
}

func RenderHtml(tmpl string) (string, error) {
	t := template.Must(template.New("tmpl").
		Funcs(AttachTemplateFuncs(template.FuncMap{}, &HtmlMarkupper{})).Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, Context{})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func HtmlColours() map[string]string {
	return map[string]string{
		"black":   "rgb(0,0,1)",
		"red":     "rgb(187,0,0)",
		"green":   "rgb(0,187,0)",
		"yellow":  "rgb(187,187,0)",
		"blue":    "rgb(0,0,187)",
		"magenta": "rgb(187,0,187)",
		"cyan":    "rgb(0,187,187)",
		"gray":    "rgb(100,100,100)",
	}
}
