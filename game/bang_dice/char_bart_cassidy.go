package bang_dice

type CharBartCassidy struct{}

func (c CharBartCassidy) Name() string {
	return "Bart Cassidy"
}

func (c CharBartCassidy) Description() string {
	return "You may take an arrow instead of losing a life point (except to Indians or Dynamite)."
}

func (c CharBartCassidy) StartingLife() int {
	return 8
}
