package starship_catan

import "testing"

func TestSectorBaseCards(t *testing.T) {
	if len(SectorBaseCards()) != 40 {
		t.Fatal("There aren't 40 sector base cards.")
	}
}

func TestSector1Cards(t *testing.T) {
	if len(Sector1Cards()) != 7 {
		t.Fatal("There aren't 7 sector 1 cards.")
	}
}

func TestSector2Cards(t *testing.T) {
	if len(Sector2Cards()) != 7 {
		t.Fatal("There aren't 7 sector 2 cards.")
	}
}

func TestSector3Cards(t *testing.T) {
	if len(Sector3Cards()) != 7 {
		t.Fatal("There aren't 7 sector 3 cards.")
	}
}

func TestSector4Cards(t *testing.T) {
	if len(Sector4Cards()) != 7 {
		t.Fatal("There aren't 7 sector 4 cards.")
	}
}

func TestShuffledSectorCards(t *testing.T) {
	if len(ShuffledSectorCards()) != 68 {
		t.Fatal("There aren't 68 sector cards in the entire deck.")
	}
}

func TestAdventure1Cards(t *testing.T) {
	cards := Adventure1Cards()
	if len(cards) != 3 {
		t.Fatal("There aren't 3 adventure 1 cards.")
	}
	for _, c := range cards {
		if _, ok := c.(Adventurer); !ok {
			t.Fatalf("Adventure card is not an Adventurer:\n\n%#v", c)
		}
	}
}

func TestAdventure2Cards(t *testing.T) {
	cards := Adventure2Cards()
	if len(cards) != 3 {
		t.Fatal("There aren't 3 adventure 2 cards.")
	}
	for _, c := range cards {
		if _, ok := c.(Adventurer); !ok {
			t.Fatalf("Adventure card is not an Adventurer:\n\n%#v", c)
		}
	}
}
