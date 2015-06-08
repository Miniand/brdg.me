package cathedral

const (
	NoPlayer        = -1
	PlayerCathedral = 2
)

type Tile struct {
	PlayerType
	Owner int
	Text  string
}

var EmptyTile = Tile{
	PlayerType{NoPlayer, 0},
	NoPlayer,
	"",
}
