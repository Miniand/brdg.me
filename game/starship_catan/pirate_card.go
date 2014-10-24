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

func (c PirateCard) FamePoints() int {
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

func (c PirateCard) Commands() []command.Command {
	return []command.Command{
		FightCommand{},
		PayCommand{},
	}
}
