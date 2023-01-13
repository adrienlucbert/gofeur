package simulation

import (
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
)

type truck struct {
	name     string
	pos      pkg.Vector
	capacity uint
	load     uint
	// TODO: away time remaining
	// TODO: state [WAITING|GONE]
}

func newTruckFromParsing(from *parsing.Truck) truck {
	return truck{
		name:     from.Name,
		pos:      pkg.Vector{X: int(from.X), Y: int(from.Y)},
		capacity: from.Weight,
		load:     0,
	}
}

func (t *truck) simulateRound(simulation *Simulation) {
}
