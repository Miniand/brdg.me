package bang_dice

type CharPaulRegret struct{}

func (c CharPaulRegret) Name() string {
	return "Paul Regret"
}

func (c CharPaulRegret) Description() string {
	return "You never lose life points to the Gatling Gun."
}

func (c CharPaulRegret) StartingLife() int {
	return 9
}
