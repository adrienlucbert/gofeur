package simulation

import (
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
)

// ParcelStatus represents the parcel's possible states
type ParcelStatus int

const (
	// StandingBy is the parcel's state when it's not carried nor dropped off
	StandingBy ParcelStatus = iota
	// Targeted is the parcel's state when it's been targeted by a forklift
	Targeted
	// Carried is the parcel's state when it's being carried by a forklift
	Carried
	// DroppedOff is the parcel's state when it's been dropped off in a truck
	DroppedOff
)

type parcel struct {
	name   string
	pos    pkg.Vector
	weight uint
	status ParcelStatus
}

// Implement prop.Pos()
func (p *parcel) Pos() pkg.Vector {
	return p.pos
}

// Implement prop.IsAvailable()
func (p *parcel) IsAvailable() bool {
	return p.status == StandingBy
}

func newParcelFromParsing(from *parsing.Parcel) parcel {
	return parcel{
		name: from.Name,
		pos:  pkg.Vector{X: int(from.X), Y: int(from.Y)},
		weight: map[parsing.Color]uint{
			parsing.Yellow: 100,
			parsing.Green:  200,
			parsing.Blue:   500,
		}[from.Color],
	}
}
