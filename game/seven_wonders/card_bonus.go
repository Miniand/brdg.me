package seven_wonders

type CardBonus struct {
	Card
	TargetKinds []int
	Directions  []int
	VP          int
	Coins       int
}

func NewCardBonus(
	name string,
	kind int,
	cost Cost,
	targetKinds, directions []int,
	vp, coins int,
	freeWith, makesFree []string,
) CardBonus {
	if targetKinds == nil || len(targetKinds) == 0 {
		panic("no targetKinds")
	}
	if directions == nil || len(directions) == 0 {
		panic("no directions")
	}
	return CardBonus{
		NewCard(name, kind, cost, freeWith, makesFree),
		targetKinds,
		directions,
		vp,
		coins,
	}
}
