package starship_catan

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
)

type GainCommand struct{}

func (c GainCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("gain", 1, input)
}

func (c GainCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanGain(p)
}

func (c GainCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	r, err := ParseResource(a[0])
	if err != nil {
		return "", err
	}
	found := false
	for _, gr := range g.GainResources {
		if gr == r {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf(
			`You aren't able to gain {{b}}%s{{_b}} at the moment`,
			ResourceNames[r])
	}
	g.GainResource(p, r)
	g.Gained(p)
	return "", nil
}

func (c GainCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	resources := make([]string, len(g.GainResources))
	for i, r := range g.GainResources {
		resources[i] = fmt.Sprintf("{{b}}%s{{_b}}", ResourceNames[r])
	}
	return fmt.Sprintf(
		"{{b}}gain ##{{_b}} to gain a resource.  Enter as much of the resource name as needed to uniquely identify it.  Eg. {{b}}gain sci{{_b}}\nYou can gain: %s",
		strings.Join(resources, ", "))
}
