package acquire

import (
	"testing"
)

func TestStart(t *testing.T) {
	g := Game{}
	if err := g.Start([]string{"Mick", "Steve"}); err != nil {
		t.Fatal(err)
	}
}

func checkCorpValues(corp int, expected map[int]int, t *testing.T) {
	for size, expectedValue := range expected {
		actual := CorpValue(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorpValue(t *testing.T) {
	low := map[int]int{
		2:  200,
		3:  300,
		4:  400,
		5:  500,
		6:  600,
		10: 600,
		11: 700,
		20: 700,
		21: 800,
		30: 800,
		31: 900,
		40: 900,
		41: 1000,
	}
	med := map[int]int{
		2:  300,
		3:  400,
		4:  500,
		5:  600,
		6:  700,
		10: 700,
		11: 800,
		20: 800,
		21: 900,
		30: 900,
		31: 1000,
		40: 1000,
		41: 1100,
	}
	high := map[int]int{
		2:  400,
		3:  500,
		4:  600,
		5:  700,
		6:  800,
		10: 800,
		11: 900,
		20: 900,
		21: 1000,
		30: 1000,
		31: 1100,
		40: 1100,
		41: 1200,
	}
	checkCorpValues(TILE_CORP_WORLDWIDE, low, t)
	checkCorpValues(TILE_CORP_SACKSON, low, t)
	checkCorpValues(TILE_CORP_FESTIVAL, med, t)
	checkCorpValues(TILE_CORP_IMPERIAL, med, t)
	checkCorpValues(TILE_CORP_AMERICAN, med, t)
	checkCorpValues(TILE_CORP_CONTINENTAL, high, t)
	checkCorpValues(TILE_CORP_TOWER, high, t)
}

func checkCorpMajorityShareholderBonuses(corp int, expected map[int]int,
	t *testing.T) {
	for size, expectedValue := range expected {
		actual := CorpMajorityShareholderBonus(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorpMajorityShareholderBonuses(t *testing.T) {
	low := map[int]int{
		2:  2000,
		3:  3000,
		4:  4000,
		5:  5000,
		6:  6000,
		10: 6000,
		11: 7000,
		20: 7000,
		21: 8000,
		30: 8000,
		31: 9000,
		40: 9000,
		41: 10000,
	}
	med := map[int]int{
		2:  3000,
		3:  4000,
		4:  5000,
		5:  6000,
		6:  7000,
		10: 7000,
		11: 8000,
		20: 8000,
		21: 9000,
		30: 9000,
		31: 10000,
		40: 10000,
		41: 11000,
	}
	high := map[int]int{
		2:  4000,
		3:  5000,
		4:  6000,
		5:  7000,
		6:  8000,
		10: 8000,
		11: 9000,
		20: 9000,
		21: 10000,
		30: 10000,
		31: 11000,
		40: 11000,
		41: 12000,
	}
	checkCorpMajorityShareholderBonuses(TILE_CORP_WORLDWIDE, low, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_SACKSON, low, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_FESTIVAL, med, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_IMPERIAL, med, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_AMERICAN, med, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_CONTINENTAL, high, t)
	checkCorpMajorityShareholderBonuses(TILE_CORP_TOWER, high, t)
}

func checkCorpMinorityShareholderBonuses(corp int, expected map[int]int,
	t *testing.T) {
	for size, expectedValue := range expected {
		actual := CorpMinorityShareholderBonus(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorpMinorityShareholderBonuses(t *testing.T) {
	low := map[int]int{
		2:  1000,
		3:  1500,
		4:  2000,
		5:  2500,
		6:  3000,
		10: 3000,
		11: 3500,
		20: 3500,
		21: 4000,
		30: 4000,
		31: 4500,
		40: 4500,
		41: 5000,
	}
	med := map[int]int{
		2:  1500,
		3:  2000,
		4:  2500,
		5:  3000,
		6:  3500,
		10: 3500,
		11: 4000,
		20: 4000,
		21: 4500,
		30: 4500,
		31: 5000,
		40: 5000,
		41: 5500,
	}
	high := map[int]int{
		2:  2000,
		3:  2500,
		4:  3000,
		5:  3500,
		6:  4000,
		10: 4000,
		11: 4500,
		20: 4500,
		21: 5000,
		30: 5000,
		31: 5500,
		40: 5500,
		41: 6000,
	}
	checkCorpMinorityShareholderBonuses(TILE_CORP_WORLDWIDE, low, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_SACKSON, low, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_FESTIVAL, med, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_IMPERIAL, med, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_AMERICAN, med, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_CONTINENTAL, high, t)
	checkCorpMinorityShareholderBonuses(TILE_CORP_TOWER, high, t)
}
