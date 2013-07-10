package lost_cities

import ()

type Game struct {
	Players         []string
	CurrentlyMoving string
	StartPlayer     string
	Board           [5][10]int
}
