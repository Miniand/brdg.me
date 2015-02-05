package seven_wonders

var AllFields = []int{
	FieldMathematics,
	FieldEngineering,
	FieldTheology,
}

type CardScience struct {
	Card
	Fields []int
}

func NewCardScience(
	name string,
	cost Cost,
	field int,
	freeWith, makesFree []string,
) CardScience {
	return CardScience{
		NewCard(name, CardKindScientific, cost, freeWith, makesFree),
		[]int{field},
	}
}

func (c CardScience) SuppString() string {
	return RenderResourceList(c.Fields, "/")
}
