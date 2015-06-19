package render

import (
	"bytes"
	"text/template"
)

type TerminalMarkupper struct {
	Markupper
	ColourStack []string
	BoldStack   int
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
	return t.Current()
}
func (t *TerminalMarkupper) StartBold() interface{} {
	t.BoldStack += 1
	return t.Current()
}
func (t *TerminalMarkupper) EndBold() interface{} {
	t.BoldStack -= 1
	return t.Current()
}
func (t *TerminalMarkupper) Current() string {
	c := "\x1b[0"
	if len(t.ColourStack) > 0 {
		c = TerminalColours()[t.ColourStack[len(t.ColourStack)-1]]
	}
	if t.BoldStack > 0 {
		c += ";1"
	} else {
		c += ";0"
	}
	return c + "m"
}

func RenderTerminal(tmpl string) (string, error) {
	t := template.Must(template.New("tmpl").
		Funcs(AttachTemplateFuncs(template.FuncMap{}, &TerminalMarkupper{})).Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, Context{})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TerminalColours() map[string]string {
	return map[string]string{
		Black:   "\x1b[30",
		Red:     "\x1b[31",
		Green:   "\x1b[32",
		Yellow:  "\x1b[33",
		Blue:    "\x1b[34",
		Magenta: "\x1b[35",
		Cyan:    "\x1b[36",
		Gray:    "\x1b[37",
	}
}
