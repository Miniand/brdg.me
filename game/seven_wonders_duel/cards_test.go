package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCards(t *testing.T) {
	for cId, c := range Cards {
		assert.NotEmpty(t, c.Name, fmt.Sprintf("card %d is missing a name", cId))
		assert.NotEmpty(t, c.Id, fmt.Sprintf("card %s is missing an Id", c.Name))
		assert.Equal(t, cId, c.Id, fmt.Sprintf("card %s doesn't have an Id matching the card key (%d != %d)", c.Name, cId, c.Id))
		assert.NotEmpty(t, c.Type, fmt.Sprintf("card %s is missing a Type", c.Name))
		if c.VPFunc != nil || c.AfterBuild != nil {
			assert.NotEmpty(t, c.Summary, fmt.Sprintf("card %s has VPFunc or AfterBuild, so it must also specify Summary", c.Name))
		}
		switch c.Type {
		case CardTypeRaw:
			assert.NotEmpty(t, c.Provides, fmt.Sprintf("card %s is of Type raw and should have Provides specified", c.Name))
			for _, goodSet := range c.Provides {
				for good, _ := range goodSet {
					if good != GoodWood &&
						good != GoodStone &&
						good != GoodClay {
						t.Errorf("card %s is of Type raw so should only provide wood, stone or clay", c.Name)
					}
				}
			}
		case CardTypeManufactured:
			assert.NotEmpty(t, c.Provides, fmt.Sprintf("card %s is of Type manufactured and should have Provides specified", c.Name))
			for _, goodSet := range c.Provides {
				for good, _ := range goodSet {
					if good != GoodGlass &&
						good != GoodPapyrus {
						t.Errorf("card %s is of Type manufactured so should only provide glass or papyrus", c.Name)
					}
				}
			}
		case CardTypeCivilian:
			assert.NotEmpty(t, c.VPRaw, fmt.Sprintf("card %s is of Type civicial and should have VPRaw specified", c.Name))
		case CardTypeScientific:
			assert.NotEmpty(t, c.Science, fmt.Sprintf("card %s is of Type scientific and should have Science specified", c.Name))
		case CardTypeCommercial:
		case CardTypeMilitary:
			assert.NotEmpty(t, c.Military, fmt.Sprintf("card %s is of Type military and should have Military specified", c.Name))
		case CardTypeGuild:
		case CardTypeWonder:
		default:
			t.Errorf("the Type of card %s is unknown", c.Name)
		}
	}
}
