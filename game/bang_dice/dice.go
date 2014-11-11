package bang_dice

const (
	DieArrow = iota
	DieDynamite
	Die1
	Die2
	DieBeer
	DieGatling
)

var DieStrings = map[int]string{
	DieArrow:    `{{b}}[{{c "green"}}A{{_c}}]{{_b}}`,
	DieDynamite: `{{b}}[{{c "red"}}D{{_c}}]{{_b}}`,
	Die1:        `{{b}}[1]{{_b}}`,
	Die2:        `{{b}}[2]{{_b}}`,
	DieBeer:     `{{b}}[{{c "yellow"}}B{{_c}}]{{_b}}`,
	DieGatling:  `{{b}}[{{c "blue"}}G{{_c}}]{{_b}}`,
}
