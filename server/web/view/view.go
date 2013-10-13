package view

import (
	"html/template"
)

var tmpl *template.Template

func init() {
	tmpl = template.New("")
	// Layout stuff
	Parse("header", headerTmpl)
	Parse("footer", footerTmpl)
	Parse("title", titleTmpl)
}

func Parse(name string, text string) *template.Template {
	if tmpl.Lookup(name) == nil {
		var err error
		if tmpl, err = tmpl.New(name).Parse(text); err != nil {
			panic(err.Error())
		}
	}
	return tmpl
}
