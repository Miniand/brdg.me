package render

import (
	"bytes"
	"fmt"
	"html/template"
)

type RgbColour struct {
	R, G, B int
}

func (c RgbColour) String() string {
	return fmt.Sprintf("%d,%d,%d", c.R, c.G, c.B)
}

func (c RgbColour) Lighten(by float64) RgbColour {
	return RgbColour{
		Lighten(c.R, by),
		Lighten(c.G, by),
		Lighten(c.B, by),
	}
}

func Lighten(val int, by float64) int {
	return val + int(float64(255-val)*by)
}

var HtmlRgbColours = map[string]RgbColour{
	Black:   RgbColour{117, 117, 117},
	Red:     RgbColour{244, 67, 54},
	Green:   RgbColour{76, 175, 80},
	Yellow:  RgbColour{249, 168, 37},
	Blue:    RgbColour{25, 118, 210},
	Magenta: RgbColour{156, 39, 176},
	Cyan:    RgbColour{0, 188, 212},
	Gray:    RgbColour{180, 180, 180},
	White:   RgbColour{255, 255, 255},
}

type HtmlMarkupper struct {
	Markupper
}

func (t *HtmlMarkupper) StartColour(colour string) interface{} {
	return template.HTML(fmt.Sprintf(
		`<span style="color:rgb(%s);">`,
		HtmlRgbColours[colour].String(),
	))
}
func (t *HtmlMarkupper) EndColour() interface{} {
	return template.HTML("</span>")
}
func (t *HtmlMarkupper) StartBg(colour string) interface{} {
	return template.HTML(fmt.Sprintf(
		`<span style="background-color:rgb(%s);">`,
		HtmlRgbColours[colour].String(),
	))
}
func (t *HtmlMarkupper) EndBg() interface{} {
	return template.HTML("</span>")
}
func (t *HtmlMarkupper) StartBold() interface{} {
	return template.HTML("<strong>")
}
func (t *HtmlMarkupper) EndBold() interface{} {
	return template.HTML("</strong>")
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
