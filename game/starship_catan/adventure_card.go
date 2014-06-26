package starship_catan

const (
	AdventurePlanetHades    = "Hades"
	AdventurePlanetPallas   = "Pallas"
	AdventurePlanetPicasso  = "Picasso"
	AdventurePlanetPoseidon = "Poseidon"
)

type AdventurePlanetCard struct {
	UnsortableCard
	Name string
}

type Adventurer interface {
	Planet() string
	Text() string
	Complete(player int, game *Game) error
}

type AdventureEnvironmentalCrisis struct {
	UnsortableCard
}

func (c AdventureEnvironmentalCrisis) Planet() string {
	return AdventurePlanetPoseidon
}

func (c AdventureEnvironmentalCrisis) Text() string {
	return "In Poseidon there are environmental problems.  Donate 1 science point and gain 3 astro and 1 resource of your choice."
}

func (c AdventureEnvironmentalCrisis) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureDiplomaticGift struct {
	UnsortableCard
}

func (c AdventureDiplomaticGift) Planet() string {
	return AdventurePlanetPicasso
}

func (c AdventureDiplomaticGift) Text() string {
	return "Greetings, Catanian!  A diplomatic gift is waiting on the planet of Picasso for you.  Gain 1 resource of your choice."
}

func (c AdventureDiplomaticGift) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureMerchantGift struct {
	UnsortableCard
}

func (c AdventureMerchantGift) Planet() string {
	return AdventurePlanetPallas
}

func (c AdventureMerchantGift) Text() string {
	return "Greetings, Catanian!  A merchant gift is waiting on the planet of Pallas for you.  Gain 1 resource of your choice."
}

func (c AdventureMerchantGift) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureFamine struct {
	UnsortableCard
}

func (c AdventureFamine) Planet() string {
	return AdventurePlanetPicasso
}

func (c AdventureFamine) Text() string {
	return "Famine on Picasso!  Donate 1 food and gain a medal and 1 resource of your choice."
}

func (c AdventureFamine) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureWholesaleOrder1 struct {
	UnsortableCard
}

func (c AdventureWholesaleOrder1) Planet() string {
	return AdventurePlanetPallas
}

func (c AdventureWholesaleOrder1) Text() string {
	return "Pallas urgently requires merchandise.  Donate 1 trade good and gain a medal and 1 resource of your choice."
}

func (c AdventureWholesaleOrder1) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventurePirateNest struct {
	UnsortableCard
}

func (c AdventurePirateNest) Planet() string {
	return AdventurePlanetHades
}

func (c AdventurePirateNest) Text() string {
	return "Pirates have take root in Hades.  Reach Hades with 4 boosters and gain a medal and 1 resource of your choice."
}

func (c AdventurePirateNest) Complete(player int, game *Game) error {
	panic("not implemented")
}
