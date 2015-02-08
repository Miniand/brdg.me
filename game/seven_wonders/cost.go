package seven_wonders

type Cost map[int]int

// Clone returns a clone of the cost.
func (c Cost) Clone() Cost {
	nc := Cost{}
	for k, v := range c {
		nc[k] = v
	}
	return nc
}

// Add adds one cost to another.
func (c Cost) Add(oc Cost) Cost {
	nc := c.Clone()
	for k, v := range oc {
		nc[k] += v
	}
	return nc
}

// Inv returns the inverse of a cost.
func (c Cost) Inv() Cost {
	nc := Cost{}
	for k, v := range c {
		nc[k] = -v
	}
	return nc
}

// Sub subtracts one cost from another.
func (c Cost) Sub(oc Cost) Cost {
	return c.Add(oc.Inv())
}

// SignSplit breaks down a cost into positive and negative components.
func (c Cost) SignSplit() (pos, neg Cost) {
	for k, v := range c {
		switch {
		case v > 0:
			pos[k] = v
		case v < 0:
			neg[k] = v
		}
	}
	return
}

// CanAfford returns whether the cost can afford another cost.
func (c Cost) CanAfford(oc Cost) bool {
	_, neg := c.Sub(oc).SignSplit()
	return len(neg) == 0
}

// Take returns a new Cost with only the specified keys.
func (c Cost) Take(keys ...int) Cost {
	nc := Cost{}
	for k := range keys {
		nc[k] = c[k]
	}
	return nc
}

// Drop returns a new Cost with the specified keys dropped.
func (c Cost) Drop(keys ...int) Cost {
	dm := map[int]bool{}
	for k := range keys {
		dm[k] = true
	}
	nc := Cost{}
	for k, v := range c {
		if !dm[k] {
			nc[k] = v
		}
	}
	return nc
}
