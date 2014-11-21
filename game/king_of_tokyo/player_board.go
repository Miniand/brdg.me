package king_of_tokyo

type PlayerBoard struct {
	Health int
	VP     int
	Energy int
	Cards  []CardBase
	Tokens []interface{}
}

func NewPlayerBoard() *PlayerBoard {
	return &PlayerBoard{
		Health: 10,
		Cards:  []CardBase{},
		Tokens: []interface{}{},
	}
}

func (b *PlayerBoard) Things() []interface{} {
	things := []interface{}{}
	for _, c := range b.Cards {
		things = append(things, c)
	}
	for _, t := range b.Tokens {
		things = append(things, t)
	}
	return AllThings(things)
}

func (b *PlayerBoard) MaxHealth() int {
	max := 10
	for _, t := range b.Things() {
		if healthMod, ok := t.(MaxHealthModifier); ok {
			max = healthMod.ModifyMaxHealth(max)
		}
	}
	return max
}

func (b *PlayerBoard) ModifyHealth(amount int) {
	b.Health += amount
	if max := b.MaxHealth(); b.Health > max {
		b.Health = max
	}
	if b.Health < 0 {
		b.Health = 0
	}
}

func (b *PlayerBoard) ModifyVP(amount int) {
	b.VP += amount
	if b.VP < 0 {
		b.VP = 0
	}
}

func (b *PlayerBoard) ModifyEnergy(amount int) {
	b.Energy += amount
	if b.Energy < 0 {
		b.Energy = 0
	}
}

func AllThings(things []interface{}) []interface{} {
	allThings := append([]interface{}{}, things...)
	for _, t := range things {
		if withThings, ok := t.(HasThings); ok {
			allThings = append(allThings, AllThings(withThings.Things())...)
		}
	}
	return allThings
}
