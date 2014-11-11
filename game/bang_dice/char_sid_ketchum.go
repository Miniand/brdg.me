package bang_dice

type CharSidKetchum struct{}

func (c CharSidKetchum) Name() string {
	return "Sid Ketchum"
}

func (c CharSidKetchum) Description() string {
	return "At the beginning of your turn, any player of your choice gains one life point."
}

func (c CharSidKetchum) StartingLife() int {
	return 8
}
