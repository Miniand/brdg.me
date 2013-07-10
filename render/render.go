package render

import (
	"text/template"
)

type Markupper interface {
	StartColour(string) string
	EndColour() string
	StartBold() string
	EndBold() string
}

// @see http://en.wikipedia.org/wiki/ANSI_escape_code#Colours
func ValidColours() []string {
	return []string{
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"gray",
	}
}

func TemplateFuncs(m Markupper) template.FuncMap {
	return template.FuncMap{
		"c": func(colour string) string {
			found := false
			for _, validColour := range ValidColours() {
				if validColour == colour {
					found = true
					break
				}
			}
			if !found {
				panic(colour + " is not a valid colour")
			}
			return m.StartColour(colour)
		},
		"_c": func() string {
			return m.EndColour()
		},
		"b": func() string {
			return m.StartBold()
		},
		"_b": func() string {
			return m.EndBold()
		},
	}
}
