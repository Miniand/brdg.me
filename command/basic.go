package command

import (
	"fmt"
	"regexp"
	"strings"
)

const ARGUMENT_REGEXP = `\s*(\b[^\s]+\b)`

type BasicCommand struct {
	Name        string
	NumArgs     int
	CanCallFunc func(player string, context interface{}) bool
	CallFunc    func(player string, context interface{}, args []string) error
	UsageFunc   func(player string, context interface{}) string
}

func (b BasicCommand) Parse(input string) []string {
	return regexp.MustCompile(fmt.Sprintf(`(?im)\s*%s`+strings.Repeat(
		ARGUMENT_REGEXP, b.NumArgs)+`\s*$`, b.Name)).FindStringSubmatch(input)
}

func (b BasicCommand) CanCall(player string, context interface{}) bool {
	return b.CanCallFunc(player, context)
}

func (b BasicCommand) Call(player string, context interface{}, args []string) error {
	return b.CallFunc(player, context, args[1:])
}

func (b BasicCommand) Usage(player string, context interface{}) string {
	if b.UsageFunc == nil {
		return b.Name
	}
	return b.UsageFunc(player, context)
}
