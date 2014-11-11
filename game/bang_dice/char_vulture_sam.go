package bang_dice

type CharVultureSam struct{}

func (c CharVultureSam) Name() string {
	return "Vulture Sam"
}

func (c CharVultureSam) Description() string {
	return "Each time another player is eliminated, you gain two life points."
}

func (c CharVultureSam) StartingLife() int {
	return 9
}
