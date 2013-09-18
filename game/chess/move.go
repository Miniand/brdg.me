package chess

type Move struct {
	From, To       Location
	TakeAt         *Location
	SubsequentMove *Move
}
