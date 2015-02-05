package seven_wonders

const (
	DirLeft  = -1
	DirDown  = 0
	DirRight = 1
)

var (
	DirAll        = []int{DirLeft, DirDown, DirRight}
	DirNeighbours = []int{DirLeft, DirRight}
	DirSelf       = []int{DirDown}
)

type CardCommercialTrade struct {
	Card
	Directions []int
	Goods      []int
}

func NewCardCommercialTrade(
	name string,
	cost Cost,
	directions, goods []int,
	freeWith, makesFree []string,
) CardCommercialTrade {
	if directions == nil || len(directions) == 0 {
		panic("no directions")
	}
	if goods == nil || len(goods) == 0 {
		panic("no goods")
	}
	return CardCommercialTrade{
		NewCard(name, CardKindCommercial, cost, freeWith, makesFree),
		directions,
		goods,
	}
}

type CardCommercialTavern struct {
	Card
}

func NewCardCommercialTavern() CardCommercialTavern {
	return CardCommercialTavern{
		NewCard(CardTavern, CardKindCommercial, nil, nil, nil),
	}
}
