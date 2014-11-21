package king_of_tokyo

type Prioritiser interface {
	Priority() int
}

// AttackModifier allows the modification of an attack immediately after roll,
// and allows the modification of the damage and who was attacked.
type AttackModifier interface {
	ModifyAttack(game *Game, player, damage int, attacked []int) (int, []int)
}

// AttackDamageForPlayerModifier allows the modification of an attack damage
// for a specific player, after any AttackModifier is run.
type AttackDamageForPlayerModifier interface {
	ModifyAttackDamageForPlayer(game *Game, player, attacked, damage int) int
}

// DamageModifier allowed the attacked monster to modify damage.
type DamageModifier interface {
	ModifyDamage(game *Game, player, attacker, damage int) int
}

type HasThings interface {
	Things() Things
}

type CardCostModifier interface {
	ModifyCardCost(game *Game, player, cost int) int
}

type PostCardBuyHandler interface {
	PostCardBuy(game *Game, player int, card CardBase, cost int)
}

type PostAttackHandler interface {
	PostAttack(game *Game, player, damage int)
}

type ExtraReroller interface {
	ExtraReroll(game *Game, player int, extra map[int]bool) map[int]bool
}

type PreResolveDiceHandler interface {
	PreResolveDice(game *Game, player int, dice []int) []int
}

type HealthZeroHandler interface {
	HealthZero(game *Game, player, zeroPlayer int)
}

type EndTurnHandler interface {
	EndTurn(game *Game, player int)
}

type MaxHealthModifier interface {
	ModifyMaxHealth(health int) int
}

type LeaveTokyoHandler interface {
	LeaveTokyo(game *Game, location, player, enteringPlayer int)
}

type DiceCountModifier interface {
	ModifyDiceCount(game *Game, player, diceCount int) int
}

type Things []interface{}

func (t Things) Len() int {
	return len(t)
}

func (t Things) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Things) Less(i, j int) bool {
	iPriority := 0
	if p, ok := t[i].(Prioritiser); ok {
		iPriority = p.Priority()
	}
	jPriority := 0
	if p, ok := t[j].(Prioritiser); ok {
		jPriority = p.Priority()
	}
	return iPriority < jPriority
}