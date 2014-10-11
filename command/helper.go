package command

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ARGUMENT_REGEXP = `([^\S\r\n]*\b[^\s]+\b)`
)

// Parses a named command with a range of args.  If minArgs or maxArgs is < 0,
// it will be unbounded in that direction.
func ParseNamedCommandRangeArgs(name string, minArgs int, maxArgs int,
	input string) []string {
	repeater := "*"
	repeaterMin := ""
	repeaterMax := ""
	if minArgs >= 0 || maxArgs >= 0 {
		if minArgs >= 0 {
			repeaterMin = fmt.Sprintf("%d", minArgs)
		}
		if maxArgs >= 0 {
			repeaterMax = fmt.Sprintf("%d", maxArgs)
		}
		repeater = fmt.Sprintf("{%s,%s}", repeaterMin, repeaterMax)
	}
	return regexp.MustCompile(fmt.Sprintf(`(?im)\A[^\S\r\n]*%s(`+ARGUMENT_REGEXP+
		`%s)[^\S\r\n]*$`, name, repeater)).FindStringSubmatch(
		strings.TrimSpace(input))
}

// Parses a named command with a specific number of arguments.
func ParseNamedCommandNArgs(name string, numArgs int, input string) []string {
	return ParseNamedCommandRangeArgs(name, numArgs, numArgs, input)
}

// Parses a named command with any number of arguments.
func ParseNamedCommand(name string, input string) []string {
	return ParseNamedCommandNArgs(name, -1, input)
}

// Parses using provided regexp, replacing spaces with non-newline space matchers
func ParseRegexp(reg, input string) []string {
	return regexp.MustCompile(fmt.Sprintf(`(?im)^[^\S\r\n]*%s[^\S\r\n]*$`,
		regexp.MustCompile(`\bARG\b`).ReplaceAllString(
			regexp.MustCompile(`\s+?`).ReplaceAllString(reg, `[^\S\r\n]+`),
			`\b[^\s]+\b`))).
		FindStringSubmatch(input)
}

// Extracts the actual arguments from the result of a ParseNamedCommand call
func ExtractNamedCommandArgs(args []string) []string {
	return regexp.MustCompile(`\s+`).Split(strings.TrimSpace(args[1]), -1)
}
