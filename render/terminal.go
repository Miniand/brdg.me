package render

import (
	"bytes"
	"fmt"
	"text/template"
)

type TerminalMarkupper struct {
	Markupper
	FgStack   []string
	BgStack   []string
	BoldStack int
}

func (t *TerminalMarkupper) StartColour(colour string) interface{} {
	t.FgStack = append(t.FgStack, colour)
	return t.Current()
}
func (t *TerminalMarkupper) EndColour() interface{} {
	if len(t.FgStack) == 0 {
		panic("There are no colours set")
	}
	t.FgStack = t.FgStack[:len(t.FgStack)-1]
	return t.Current()
}
func (t *TerminalMarkupper) StartBg(colour string) interface{} {
	t.BgStack = append(t.BgStack, colour)
	return t.Current()
}
func (t *TerminalMarkupper) EndBg() interface{} {
	if len(t.BgStack) == 0 {
		panic("There are no colours set")
	}
	t.BgStack = t.BgStack[:len(t.BgStack)-1]
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
	fgCode := 39
	bgCode := 49
	boldCode := ""
	if len(t.FgStack) > 0 {
		fgCode = AnsiFgCode(t.FgStack[len(t.FgStack)-1])
	}
	if len(t.BgStack) > 0 {
		bgCode = AnsiBgCode(t.BgStack[len(t.BgStack)-1])
	}
	if t.BoldStack > 0 {
		boldCode = ";1"
	}
	return fmt.Sprintf("\x1b[0;%d;%d%sm", fgCode, bgCode, boldCode)
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

var AnsiColourNums = map[string]int{
	Black:   0,
	Red:     1,
	Green:   2,
	Yellow:  3,
	Blue:    4,
	Magenta: 5,
	Cyan:    6,
	Gray:    7,
	White:   9,
}

func AnsiFgCode(colour string) int {
	return 30 + AnsiColourNums[colour]
}

func AnsiBgCode(colour string) int {
	if colour == White {
		return AnsiBgCode(Gray)
	}
	return 40 + AnsiColourNums[colour]
}
