package starship_catan

import (
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type ColonyCard struct {
	UnsortableCard
	Name      string
	Resource  int
	Dice      int
	StartCard bool
}

func (c ColonyCard) String() string {
	return fmt.Sprintf(
		`{{c "green"}}{{b}}%s{{_b}}{{_c}} (colony planet, roll {{b}}%d{{_b}} for {{b}}%s{{_b}})`,
		c.Name,
		c.Dice,
		RenderResource(c.Resource),
	)
}

func (c ColonyCard) Commands() []command.Command {
	return []command.Command{
		FoundColonyCommand{},
	}
}

func (c ColonyCard) VictoryPoints() int {
	return 1
}
