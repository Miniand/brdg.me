package king_of_tokyo

type CardPlotTwist struct{}

func (c CardPlotTwist) Name() string {
	return "Plot Twist"
}

func (c CardPlotTwist) Description() string {
	return "{{b}}Change one die to any result.{{_b}} Discard when used."
}

func (c CardPlotTwist) Cost() int {
	return 3
}

func (c CardPlotTwist) Kind() int {
	return CardKindKeep
}
