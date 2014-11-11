package bang_dice

type CharPedroRamirez struct{}

func (c CharPedroRamirez) Name() string {
	return "Pedro Ramirez"
}

func (c CharPedroRamirez) Description() string {
	return "Each time you lose a life point, you may discard one of your arrows."
}

func (c CharPedroRamirez) StartingLife() int {
	return 8
}
