package starship_catan

type MedianCard struct {
	UnsortableCard
}

func (c MedianCard) CanFoundTradingPost() bool {
	return true
}

func (c MedianCard) DiplomatPoints() int {
	return 2
}
