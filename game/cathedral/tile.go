package cathedral

const (
	NoPlayer        = -1
	PlayerCathedral = 2
)

type Tile struct {
	Player int
	Type   int
}
