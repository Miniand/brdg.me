package starship_catan

import "fmt"

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
		`a {{c "gray"}}{{b}}pirate ship{{_b}}{{_c}}, asking a ransom of %s`,
		RenderMoney(c.Ransom),
	)
}
