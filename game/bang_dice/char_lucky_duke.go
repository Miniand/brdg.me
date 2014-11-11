package bang_dice

type CharLuckyDuke struct{}

func (c CharLuckyDuke) Name() string {
	return "Lucky Duke"
}

func (c CharLuckyDuke) Description() string {
	return "You may take one extra re-roll."
}

func (c CharLuckyDuke) StartingLife() int {
	return 8
}
