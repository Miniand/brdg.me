package starship_catan

import "errors"

const (
	TradeDirBoth = iota
	TradeDirBuy
	TradeDirSell
)

type TradeCard struct {
	UnsortableCard
	Name        string
	Resource    int
	Price       int
	Minimum     int
	Maximum     int
	Direction   int
	TradingPost bool
}

func (c TradeCard) Buy(resource, amount int) error {
	return errors.New("not implemented")
}

func (c TradeCard) Sell(resource, amount int) error {
	return errors.New("not implemented")
}

func (c TradeCard) FriendshipPoints() int {
	if c.TradingPost {
		return 1
	}
	return 0
}

func (c TradeCard) CanFoundTradingPost() bool {
	return c.TradingPost
}
