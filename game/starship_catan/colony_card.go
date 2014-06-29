package starship_catan

type ColonyCard struct {
	UnsortableCard
	Name      string
	Resource  int
	Dice      int
	StartCard bool
}

func (c ColonyCard) VictoryPoints() int {
	if c.StartCard {
		return 0
	}
	return 1
}
