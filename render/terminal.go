package render

import (
	"bytes"
	"github.com/beefsack/boredga.me/game"
	"text/template"
)

type TerminalMarkupper struct {
	Markupper
	ColourStack []string
	Bold        bool
}

func (t *TerminalMarkupper) StartColour(colour string) interface{} {
	t.ColourStack = append(t.ColourStack, colour)
	return t.Current()
}
func (t *TerminalMarkupper) EndColour() interface{} {
	if len(t.ColourStack) == 0 {
		panic("There are no colours set")
	}
	t.ColourStack = t.ColourStack[:len(t.ColourStack)-1]
	if len(t.ColourStack) == 0 {
		return t.Current()
	}
	return t.ColourStack[len(t.ColourStack)-1]
}
func (t *TerminalMarkupper) StartBold() interface{} {
	t.Bold = true
	return t.Current()
}
func (t *TerminalMarkupper) EndBold() interface{} {
	t.Bold = false
	return t.Current()
}
func (t *TerminalMarkupper) Current() string {
	c := "\x1b[0"
	if len(t.ColourStack) > 0 {
		c = TerminalColours()[t.ColourStack[len(t.ColourStack)-1]]
	}
	if t.Bold {
		c += ";1"
	}
	return c + "m"
}

func RenderTerminal(tmpl string, g game.Playable) (string, error) {
	t := template.Must(template.New("tmpl").
		Funcs(AttachTemplateFuncs(template.FuncMap{}, &TerminalMarkupper{})).Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, g)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TerminalColours() map[string]string {
	return map[string]string{
		"black":   "\x1b[30",
		"red":     "\x1b[31",
		"green":   "\x1b[32",
		"yellow":  "\x1b[33",
		"blue":    "\x1b[34",
		"magenta": "\x1b[35",
		"cyan":    "\x1b[36",
		"gray":    "\x1b[37",
	}
}
