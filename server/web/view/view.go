package view

import (
	"html/template"
)

var tmpl *template.Template

var funcs = template.FuncMap{
	"loggedInUser": func() interface{} {
		if LoggedInUser == "" {
			return nil
		}
		return LoggedInUser
	},
}

var LoggedInUser = ""

func init() {
	tmpl = template.New("")
	// Layout stuff
	Parse("header", headerTmpl)
	Parse("footer", footerTmpl)
	Parse("title", titleTmpl)
}

func Parse(name string, text string) *template.Template {
	t := tmpl.Lookup(name)
	if t == nil {
		t = template.Must(tmpl.New(name).Funcs(funcs).Parse(text))
	}
	return t
}
