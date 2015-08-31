package starship_catan

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
)

type PirateCard struct {
	UnsortableCard
	Strength      int
	Ransom        int
	DestroyCannon bool
	DestroyModule bool
}

func (c PirateCard) Medals() int {
	return 1
}

func (c PirateCard) String() string {
	return fmt.Sprintf(
		`{{c "gray"}}{{b}}pirate ship{{_b}}{{_c}}, asking a ransom of %s`,
		RenderMoney(c.Ransom),
	)
}

func (c PirateCard) FullString() string {
	extra := []string{
		fmt.Sprintf(`strength {{b}}%d{{_b}}`, c.Strength),
	}
	if c.DestroyCannon {
		extra = append(extra, "destroys cannon")
	}
	if c.DestroyModule {
		extra = append(extra, "destroys module")
	}
	return fmt.Sprintf(
		`%s (%s)`,
		c,
		strings.Join(extra, ", "),
	)
}

func (c PirateCard) RequiresAction() bool {
	return true
}

func (c PirateCard) Commands(g *Game, player int) []command.Command {
	commands := []command.Command{}
	if g.CanFight(player) {
		commands = append(commands, FightCommand{})
	}
	if g.CanPayRansom(player) {
		commands = append(commands, PayCommand{})
	}
	if g.CanLoseModule(player) {
		commands = append(commands, LoseCommand{})
	}
	return commands
}
