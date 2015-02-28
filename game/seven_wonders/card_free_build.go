package seven_wonders

import "encoding/gob"

type CardFreeBuild struct {
	Card
	HasBuilt bool
}

func init() {
	gob.Register(&CardFreeBuild{})
}

func (c *CardFreeBuild) CanFreeBuild() bool {
	return !c.HasBuilt
}

func (c *CardFreeBuild) HandleFreeBuild() {
	c.HasBuilt = true
}

func (c *CardFreeBuild) SuppString() string {
	return "Build for free once each round"
}

func (c *CardFreeBuild) HandleStartRound() {
	c.HasBuilt = false
}
