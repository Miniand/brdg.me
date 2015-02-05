package seven_wonders

type CardMilitary struct {
	Card
	Strength int
}

func NewCardMilitary(
	name string,
	cost Cost,
	strength int,
	freeWith, makesFree []string,
) CardMilitary {
	return CardMilitary{
		NewCard(name, CardKindMilitary, cost, freeWith, makesFree),
		strength,
	}
}
