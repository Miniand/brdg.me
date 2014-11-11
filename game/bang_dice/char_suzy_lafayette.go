package bang_dice

import "fmt"

type CharSuzyLafayette struct{}

func (c CharSuzyLafayette) Name() string {
	return "Suzy Lafayette"
}

func (c CharSuzyLafayette) Description() string {
	return fmt.Sprintf(
		"If you didn't roll any %s or %s you gain two life points.",
		DieStrings[Die1],
		DieStrings[Die2],
	)
}

func (c CharSuzyLafayette) StartingLife() int {
	return 8
}
