package view

import (
	"io"
)

var rootTmpl = `{{template "header" "Welcome"}}
<div class="title">
	<h1>{{template "title"}}

Play board games online by web and email</h1>
</div>
{{template "footer"}}`

func Root(wr io.Writer) {
	Parse("root", rootTmpl).Execute(wr, nil)
}
