package starship_catan

import "github.com/Miniand/brdg.me/command"

type MedianCard struct {
	UnsortableCard
}

func (c MedianCard) CanFoundTradingPost() bool {
	return true
}

func (c MedianCard) DiplomatPoints() int {
	return 2
}

func (c MedianCard) Commands() []command.Command {
	return []command.Command{
		FoundTradeCommand{},
	}
}
