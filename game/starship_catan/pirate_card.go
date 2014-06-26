package starship_catan

type PirateCard struct {
	UnsortableCard
	Strength      int
	Ransom        int
	DestroyCannon bool
	DestroyModule bool
}

func (c PirateCard) FamePoints() int {
	return 1
}
