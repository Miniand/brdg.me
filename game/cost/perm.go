package cost

func CanAffordPerm(c Cost, with [][]Cost) (can bool, canWith [][]Cost) {
	canWith = [][]Cost{}
	if c.IsZero() {
		can = true
		return
	}
	if with == nil || len(with) == 0 {
		return
	}
	wLen := 0
	if with[0] != nil {
		wLen = len(with[0])
	}
	switch wLen {
	case 0:
		return CanAffordPerm(c, with[1:])
	case 1:
		using := with[0][0].Take(c.Keys()...)
		if using.CanAfford(c) {
			// Can afford it, exit early
			return true, [][]Cost{{c}}
		}
		remaining, _ := c.Sub(using).PosNeg()
		used := c.Sub(remaining)
		if subCan, subCanWith := CanAffordPerm(
			remaining,
			with[1:],
		); subCan {
			can = true
			for _, scw := range subCanWith {
				canWith = append(canWith, append([]Cost{used}, scw...))
			}
		}
	default:
		canWith = [][]Cost{}
		for _, w := range with[0] {
			if subCan, subCanWith := CanAffordPerm(
				c,
				append([][]Cost{{w}}, with[1:]...),
			); subCan {
				can = true
				canWith = append(canWith, subCanWith...)
			}
		}
	}
	return
}
