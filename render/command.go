package render

import (
	"bytes"
	"github.com/beefsack/brdg.me/command"
)

func OutputCommands(player string, context interface{},
	commands []command.Command) string {
	buf := bytes.NewBufferString("{{b}}You can:{{_b}}\n")
	for _, c := range commands {
		buf.WriteString(" * ")
		buf.WriteString(c.Usage(player, context))
		buf.WriteByte('\n')
	}
	return buf.String()
}
