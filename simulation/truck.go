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
	awayTime uint
	awayLeft uint
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
		awayTime: from.RAvail,
		awayLeft: 0,
	}
}

func (t *truck) startDelivery() {
	t.status = Away
	t.awayLeft = t.awayTime
}

func (t *truck) simulateRound(simulation *Simulation) {
	switch t.status {
	case Loading:
		if t.load > 0 {
			t.startDelivery()
		}
		// TODO: determine whether or not it is profitable to start delivery, based on:
		// - distance to nearest forklift and its load
		// - distance to nearest parcel and its weight
		// availableLoad := t.capacity - t.load
		// if target := findClosestParcel(simulation.parcels, t.pos); target != nil {
		// 	if target.weight > availableLoad {
		// 		t.startDelivery()
		// 	}
		// }
	case Away:
		t.awayLeft--
		if t.awayLeft == 0 {
			t.load = 0
			t.status = Loading
		}
	}
}
