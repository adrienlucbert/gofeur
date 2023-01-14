package simulation

import (
	"github.com/adrienlucbert/gofeur/logger"
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
	name         string
	pos          pkg.Vector
	capacity     uint
	load         uint
	loadEstimate uint
	status       TruckStatus
	awayTime     uint
	awayLeft     uint
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
		name:         from.Name,
		pos:          pkg.Vector{X: int(from.X), Y: int(from.Y)},
		capacity:     uint(from.MaxWeight),
		load:         0,
		loadEstimate: 0,
		status:       Loading,
		awayTime:     uint(from.Available),
		awayLeft:     0,
	}
}

func (t *truck) startDelivery() {
	t.status = Away
	t.awayLeft = t.awayTime
}

func (t *truck) simulateRound(simulation *Simulation) {
	switch t.status {
	case Loading:
		availableLoad := t.capacity - t.loadEstimate
		var parcelIsNearby bool
		if target := findClosestParcel(simulation.parcels, t.pos, availableLoad); target != nil {
			// determine if a forklift would have roughly enough time to travel from
			// the truck to nearest parcel and back in the time the truck would be away
			parcelIsNearby = t.pos.Distance(target.pos) <= float32(t.awayTime)*2
		}
		if t.load > 0 && t.load == t.loadEstimate && !parcelIsNearby {
			t.startDelivery()
		}
	case Away:
		t.awayLeft--
		if t.awayLeft == 0 {
			t.load = 0
			t.loadEstimate = 0
			t.status = Loading
		}
	}
	action := map[TruckStatus]string{
		Loading: "WAITING",
		Away:    "GONE",
	}[t.status]
	logger.Info("%s %s %d/%d\n", t.name, action, t.load, t.capacity)
}
