package splendor

type Amount map[int]int

func (a Amount) Clone() Amount {
	newAmount := Amount{}
	for r, n := range a {
		newAmount[r] = n
	}
	return newAmount
}

func (a Amount) Gems() Amount {
	newAmount := Amount{}
	for _, g := range Gems {
		if n, ok := a[g]; ok {
			newAmount[g] = n
		}
	}
	return newAmount
}

func (a Amount) Subtract(b Amount) Amount {
	newAmount := a.Clone()
	for r, n := range b {
		newAmount[r] -= n
	}
	return newAmount
}

func (a Amount) Add(b Amount) Amount {
	newAmount := a.Clone()
	for r, n := range b {
		newAmount[r] += n
	}
	return newAmount
}

func (a Amount) CanAfford(cost Amount) bool {
	short := 0
	for g, n := range cost {
		if a[g] < n {
			short += n - a[g]
		}
	}
	return a[Gold]-cost[Gold] >= short
}
