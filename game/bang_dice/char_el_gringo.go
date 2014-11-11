package bang_dice

type CharElGringo struct{}

func (c CharElGringo) Name() string {
	return "El Gringo"
}

func (c CharElGringo) Description() string {
	return "When a player makes you lose one or more life points, he must take an arrow."
}

func (c CharElGringo) StartingLife() int {
	return 7
}
