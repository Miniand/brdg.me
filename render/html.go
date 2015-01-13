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
		Black:   "#000000",
		Red:     "#F44336",
		Green:   "#4CAF50",
		Yellow:  "#F9A825",
		Blue:    "#2196F3",
		Magenta: "#9C27B0",
		Cyan:    "#00BCD4;",
		Gray:    "#757575",
	}
}
