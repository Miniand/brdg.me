package render

import (
	"bytes"
	"github.com/Miniand/brdg.me/command"
)

func OutputCommands(player string, context interface{},
	commands []command.Command) string {
	buf := bytes.NewBufferString("{{b}}You can:{{_b}}")
	for _, c := range commands {
		buf.WriteByte('\n')
		buf.WriteString(" * ")
		buf.WriteString(c.Usage(player, context))
	}
	return buf.String()
}
