package starship_catan

type ColonyCard struct {
	UnsortableCard
	Name     string
	Resource int
	Dice     int
}

func (c ColonyCard) VictoryPoints() int {
	return 1
}
