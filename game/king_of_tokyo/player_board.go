package king_of_tokyo

import "sort"

type PlayerBoard struct {
	Health int
	VP     int
	Energy int
	Cards  []int
	Tokens []interface{}

	// Card specific state
	HasDealtDamage bool
}

func NewPlayerBoard() *PlayerBoard {
	return &PlayerBoard{
		Health: 10,
		Cards:  []int{},
		Tokens: []interface{}{},
	}
}

func (b *PlayerBoard) Things() Things {
	things := Things{}
	for _, c := range b.Cards {
		things = append(things, Cards[c])
	}
	for _, t := range b.Tokens {
		things = append(things, t)
	}
	allThings := AllThings(things)
	sort.Sort(allThings)
	return allThings
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

func AllThings(things Things) Things {
	allThings := append(Things{}, things...)
	for _, t := range things {
		if withThings, ok := t.(HasThings); ok {
			allThings = append(allThings, AllThings(withThings.Things())...)
		}
	}
	return allThings
}
