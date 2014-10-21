package starship_catan

import (
	"errors"
	"fmt"
)

const (
	TradeDirBoth = iota
	TradeDirBuy
	TradeDirSell
)

var TradeDirStrings = map[int]string{
	TradeDirBoth: "buy/sell",
	TradeDirBuy:  "buy",
	TradeDirSell: "sell",
}

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

func (c TradeCard) String() string {
	amount := ""
	switch {
	case c.Minimum > 0 && c.Minimum == c.Maximum:
		amount = fmt.Sprintf(` {{b}}%d{{_b}}`, c.Minimum)
	case c.Minimum > 0 && c.Maximum > 0:
		amount = fmt.Sprintf(
			` between {{b}}%d{{_b}} and {{b}}%d{{_b}}`,
			c.Minimum,
			c.Maximum,
		)
	case c.Minimum > 0:
		amount = fmt.Sprintf(` at least {{b}}%d{{_b}}`, c.Minimum)
	case c.Maximum > 0:
		amount = fmt.Sprintf(` up to {{b}}%d{{_b}}`, c.Maximum)
	}
	return fmt.Sprintf(
		`{{c "yellow"}}{{b}}%s{{_b}}{{_c}} (%s%s {{b}}%s{{_b}} for %s each)`,
		c.Name,
		TradeDirStrings[c.Direction],
		amount,
		RenderResource(c.Resource),
		RenderMoney(c.Price),
	)
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
