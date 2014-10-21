package starship_catan

import "fmt"

const (
	AdventurePlanetHades    = "Hades"
	AdventurePlanetPallas   = "Pallas"
	AdventurePlanetPicasso  = "Picasso"
	AdventurePlanetPoseidon = "Poseidon"
)

var AdventurePlanetColours = map[string]string{
	AdventurePlanetHades:    "red",
	AdventurePlanetPallas:   "yellow",
	AdventurePlanetPicasso:  "magenta",
	AdventurePlanetPoseidon: "cyan",
}

type AdventurePlanetCard struct {
	UnsortableCard
	Name string
}

func (c AdventurePlanetCard) String() string {
	return fmt.Sprintf(
		`{{c "%s"}}{{b}}%s{{_b}}{{_c}}`,
		AdventurePlanetColours[c.Name],
		c.Name,
	)
}

type Adventurer interface {
	Planet() string
	Text() string
	Complete(player int, game *Game) error
}

// Adventure deck 1

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

// Adventure deck 2

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

// Adventure deck 3

type AdventureCouncilMeeting struct {
	UnsortableCard
}

func (c AdventureCouncilMeeting) Planet() string {
	return AdventurePlanetPoseidon
}

func (c AdventureCouncilMeeting) Text() string {
	return "The Galactic Council urgently requires 6 Astro to organise the meeting of the council.  Donate 6 Astro and gain a medal and 2 resources of your choice."
}

func (c AdventureCouncilMeeting) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureEpidemic struct {
	UnsortableCard
}

func (c AdventureEpidemic) Planet() string {
	return AdventurePlanetHades
}

func (c AdventureEpidemic) Text() string {
	return "A mystery plague has broken out on Hades.  Donate 2 science points and gain a victory point."
}

func (c AdventureEpidemic) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureEmergency struct {
	UnsortableCard
}

func (c AdventureEmergency) Planet() string {
	return AdventurePlanetPicasso
}

func (c AdventureEmergency) Text() string {
	return "A spaceship near Picasso is in a gravitational trap.  Whoever reaches picasso with 4 boosters can set them free and gain a medal and 1 resource of your choice."
}

func (c AdventureEmergency) Complete(player int, game *Game) error {
	panic("not implemented")
}

// Adventure deck 4

type AdventureReconstruction struct {
	UnsortableCard
}

func (c AdventureReconstruction) Planet() string {
	return AdventurePlanetHades
}

func (c AdventureReconstruction) Text() string {
	return "We have freed Hades from pirates and the population urgently requires reconstruction aid.  Donate 10 Astro and gain 2 medals."
}

func (c AdventureReconstruction) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureMonument struct {
	UnsortableCard
}

func (c AdventureMonument) Planet() string {
	return AdventurePlanetPallas
}

func (c AdventureMonument) Text() string {
	return "The Pallas population wants to build a monument for the merchants.  Donate 2 ore and 1 carbon and gain a victory point."
}

func (c AdventureMonument) Complete(player int, game *Game) error {
	panic("not implemented")
}

type AdventureWholesaleOrder2 struct {
	UnsortableCard
}

func (c AdventureWholesaleOrder2) Planet() string {
	return AdventurePlanetPoseidon
}

func (c AdventureWholesaleOrder2) Text() string {
	return "This time Pallas urgently requires merchandise.  Donate 2 trade goods and gain a medal and 2 resources of your choice."
}

func (c AdventureWholesaleOrder2) Complete(player int, game *Game) error {
	panic("not implemented")
}
