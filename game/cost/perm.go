package cost

func prependToCostArrays(c Cost, arr [][]Cost) [][]Cost {
	l := len(arr)
	ret := make([][]Cost, l)
	for i, a := range arr {
		ret[i] = append([]Cost{c}, a...)
	}
	return ret
}

func CanAffordPerm(c Cost, with [][]Cost) (can bool, canWith [][]Cost) {
	canWith = [][]Cost{}
	if c.IsZero() {
		can = true
		return
	}
	if len(with) == 0 {
		return
	}

	canWith = [][]Cost{}
	relevant := false
	cKeys := c.Keys()
	for _, w := range with[0] {
		if w.CanAfford(c) {
			return true, [][]Cost{{c}}
		}
		needed := w.Take(cKeys...).Trim()
		if len(needed) == 0 {
			continue
		}
		relevant = true
		remaining, extra := c.Sub(needed).PosNeg()
		if subCan, subCanWith := CanAffordPerm(
			remaining,
			with[1:],
		); subCan {
			can = true
			canWith = append(canWith, prependToCostArrays(
				needed.Add(extra),
				subCanWith,
			)...)
		}
	}
	if !relevant {
		subCan, subCanWith := CanAffordPerm(c, with[1:])
		can = subCan
		canWith = prependToCostArrays(Cost{}, subCanWith)
	}
	return
}
