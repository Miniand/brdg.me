package starship_catan

import (
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

const (
	TradeDirBoth = 0
	TradeDirBuy  = 1
	TradeDirSell = -1
)

var TradeDirStrings = map[int]string{
	TradeDirBoth: "buy/sell",
	TradeDirBuy:  "buy",
	TradeDirSell: "sell",
}

var TradeDirPastStrings = map[int]string{
	TradeDirBoth: "bought/sold",
	TradeDirBuy:  "bought",
	TradeDirSell: "sold",
}

func AmountTradeDir(amount int) int {
	switch {
	case amount == 0:
		return TradeDirBoth
	case amount > 0:
		return TradeDirBuy
	default:
		return TradeDirSell
	}
}

type TradeCard struct {
	UnsortableCard
	Name        string
	Resources   []int
	Price       int
	Maximum     int
	Direction   int
	TradingPost bool
}

func (c TradeCard) AmountLimitString() string {
	switch {
	case c.Maximum > 0:
		return fmt.Sprintf(`up to {{b}}%d{{_b}}`, c.Maximum)
	default:
		return ""
	}
}

func (c TradeCard) String() string {
	amount := ""
	if c.Maximum > 0 {
		amount = fmt.Sprintf(" %s", c.AmountLimitString())
	}
	return fmt.Sprintf(
		`{{c "yellow"}}{{b}}%s{{_b}}{{_c}} (%s%s {{b}}%s{{_b}} for %s each)`,
		c.Name,
		TradeDirStrings[c.Direction],
		amount,
		RenderResources(c.Resources),
		RenderMoney(c.Price),
	)
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

func (c TradeCard) Commands() []command.Command {
	return []command.Command{
		BuyCommand{},
		SellCommand{},
		NextCommand{},
	}
}
