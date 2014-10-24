package starship_catan

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/log"
)

type Transaction map[int]int

// Inverse gets the inverse of the transaction.
func (t Transaction) Inverse() Transaction {
	for k, v := range t {
		t[k] = v * -1
	}
	return t
}

func (t Transaction) String() string {
	gStr := t.GainString()
	lStr := t.LoseString()
	switch {
	case gStr != "" && lStr != "":
		return fmt.Sprintf("got %s for %s", gStr, lStr)
	case gStr != "":
		return fmt.Sprintf("got %s", gStr)
	case lStr != "":
		return fmt.Sprintf("paid %s", lStr)
	}
	return ""
}

func (t Transaction) GainString() string {
	gain := []string{}
	for _, r := range Resources {
		v := t[r]
		if v > 0 {
			gain = append(gain, RenderResourceAmount(r, v))
		}
	}
	return strings.Join(gain, ", ")
}

func (t Transaction) LoseString() string {
	lose := []string{}
	for _, r := range Resources {
		v := t[r]
		if v < 0 {
			lose = append(lose, RenderResourceAmount(r, -v))
		}
	}
	return strings.Join(lose, ", ")
}

func (t Transaction) IsEmpty() bool {
	for _, v := range t {
		if v != 0 {
			return false
		}
	}
	return true
}

func (t Transaction) TrimEmpty() Transaction {
	trimmed := Transaction{}
	for r, v := range t {
		if v != 0 {
			trimmed[r] = v
		}
	}
	return trimmed
}

func (t Transaction) Resources() []int {
	resources := []int{}
	for r, _ := range t.TrimEmpty() {
		resources = append(resources, r)
	}
	return resources
}

func (t Transaction) Gain() Transaction {
	gain := Transaction{}
	for r, v := range t {
		if v > 0 {
			gain[r] = v
		}
	}
	return gain
}

func (t Transaction) Lose() Transaction {
	lose := Transaction{}
	for r, v := range t {
		if v < 0 {
			lose[r] = v
		}
	}
	return t
}

func TransactionFromResources(resources []int) Transaction {
	t := Transaction{}
	for _, r := range resources {
		t[r] = 1
	}
	return t
}

func (t Transaction) CannotAffordError() error {
	return fmt.Errorf("can't afford %s", t.LoseString())
}

func (t Transaction) CannotFitError() error {
	return fmt.Errorf("not enough room for %s", t.GainString())
}

func (g *Game) LogTransaction(player int, t Transaction) {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s %s`,
		g.RenderName(player),
		t,
	)))
}
