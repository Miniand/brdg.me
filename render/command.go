package render

import (
	"bytes"
)

func CommandUsages(usages []string) string {
	buf := bytes.NewBufferString("{{b}}You can:{{_b}}")
	for _, usage := range usages {
		buf.WriteByte('\n')
		buf.WriteString(" * ")
		buf.WriteString(usage)
	}
	return buf.String()
}
