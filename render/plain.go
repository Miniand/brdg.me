package render

import "regexp"

var plainReplaceRegexp = regexp.MustCompile(`\{\{[^}]*\}\}`)

func RenderPlain(tmpl string) string {
	return plainReplaceRegexp.ReplaceAllString(tmpl, "")
}
