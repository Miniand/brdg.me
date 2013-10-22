package render

import (
	"fmt"
	"strings"
)

func CommandUsages(usages []string) string {
	rendered := []string{}
	for _, usage := range usages {
		rendered = append(rendered, fmt.Sprintf(" * %s", usage))
	}
	return strings.Join(rendered, "\n")
}
