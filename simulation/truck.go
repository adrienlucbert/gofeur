package simulation

import (
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
)

// TruckStatus represents the truck's possible states
type TruckStatus int

const (
	// Loading is the truck's state when it's idle, loading parcels
	Loading TruckStatus = iota
	// Away is the truck's state when it's away, delivering parcels
	Away
)

type truck struct {
	name     string
	pos      pkg.Vector
	capacity uint
	load     uint
	status   TruckStatus
	// TODO: away time remaining
}

// Implement prop.Pos()
func (t *truck) Pos() pkg.Vector {
	return t.pos
}

// Implement prop.IsAvailable()
func (t *truck) IsAvailable() bool {
	return t.status == Loading
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
