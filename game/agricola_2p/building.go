package agricola_2p

import "encoding/gob"

func init() {
	gob.Register(Cottage{})
}

type Building interface {
	Capacity() int
	String() string
}

type Cottage struct{}

func (c Cottage) Capacity() int {
	return 1
}

func (c Cottage) String() string {
	return "Cottage"
}
