package cost

import "sort"

type Cost map[int]int

// Create a cost from a slice of ints.
func FromInts(ints []int) Cost {
	c := Cost{}
	for _, i := range ints {
		c[i]++
	}
	return c
}

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

// PosNeg breaks down a cost into positive and negative components.
func (c Cost) PosNeg() (pos, neg Cost) {
	pos, neg = Cost{}, Cost{}
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
	_, neg := c.Sub(oc).PosNeg()
	return len(neg) == 0
}

// Take returns a new Cost with only the specified keys.
func (c Cost) Take(keys ...int) Cost {
	nc := Cost{}
	for _, k := range keys {
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

// Keys gets all keys for non-zero values.
func (c Cost) Keys() []int {
	keys := []int{}
	for k, v := range c {
		if v != 0 {
			keys = append(keys, k)
		}
	}
	sort.Ints(keys)
	return keys
}

// IsZero returns whether this cost only contains zero values.
func (c Cost) IsZero() bool {
	for _, v := range c {
		if v != 0 {
			return false
		}
	}
	return true
}

// Trim removes any zero values.
func (c Cost) Trim() Cost {
	nc := Cost{}
	for k, v := range c {
		if v != 0 {
			nc[k] = v
		}
	}
	return nc
}

// Ints returns an int slice of the cost keys, repeated based on the value and
// sorted.
func (c Cost) Ints() []int {
	ints := []int{}
	for k, v := range c {
		for i := 0; i < v; i++ {
			ints = append(ints, k)
		}
	}
	sort.Ints(ints)
	return ints
}

// Sum is a sum of all of the values in the cost.
func (c Cost) Sum() int {
	sum := 0
	for _, v := range c {
		sum += v
	}
	return sum
}

// Mul multiplies the cost by a number.
func (c Cost) Mul(by int) Cost {
	nc := c.Clone()
	for k := range c {
		nc[k] *= by
	}
	return nc
}
