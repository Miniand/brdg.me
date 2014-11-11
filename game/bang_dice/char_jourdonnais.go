package bang_dice

type CharJourdonnais struct{}

func (c CharJourdonnais) Name() string {
	return "Jourdonnais"
}

func (c CharJourdonnais) Description() string {
	return "You never lose more than one life point to Indians."
}

func (c CharJourdonnais) StartingLife() int {
	return 7
}
