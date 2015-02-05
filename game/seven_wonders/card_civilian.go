package seven_wonders

type CardCivilian struct {
	Card
	VP int
}

func NewCardCivilian(
	name string,
	cost Cost,
	vp int,
	freeWith, makesFree []string,
) CardCivilian {
	return CardCivilian{
		NewCard(name, CardKindCivilian, cost, freeWith, makesFree),
		vp,
	}
}
