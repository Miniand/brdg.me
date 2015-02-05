package seven_wonders

type CardGood struct {
	Card
	Goods  []int
	Amount int
}

func NewCardGood(
	name string,
	kind int,
	cost Cost,
	goods []int,
	amount int,
	freeWith, makesFree []string,
) CardGood {
	if goods == nil || len(goods) == 0 {
		panic("no goods")
	}
	return CardGood{
		NewCard(name, kind, cost, freeWith, makesFree),
		goods,
		amount,
	}
}

func NewCardGoodRaw(
	name string,
	cost Cost,
	goods []int,
	amount int,
) CardGood {
	return NewCardGood(
		name,
		CardKindRaw,
		cost,
		goods,
		amount,
		nil,
		nil,
	)
}

func NewCardGoodManufactured(
	name string,
	cost Cost,
	goods []int,
	amount int,
) CardGood {
	return NewCardGood(
		name,
		CardKindManufactured,
		cost,
		goods,
		amount,
		nil,
		nil,
	)
}

func NewCardGoodCommercial(
	name string,
	cost Cost,
	goods []int,
	amount int,
	freeWith, makesFree []string,
) CardGood {
	return NewCardGood(
		name,
		CardKindCommercial,
		cost,
		goods,
		amount,
		freeWith,
		makesFree,
	)
}
