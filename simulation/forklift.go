package simulation

import (
	"fmt"

	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pathfinding"
	"github.com/adrienlucbert/gofeur/pkg"
)

type forklift struct {
	name string
	pos  pkg.Vector
}

func newForkliftFromParsing(from *parsing.Forklift) forklift {
	return forklift{
		name: from.Name,
		pos:  pkg.Vector{X: int(from.X), Y: int(from.Y)},
	}
}

func (f *forklift) simulateRound(simulation *Simulation) {
	closestParcel := findClosestParcel(simulation.parcels, f.pos)
	if closestParcel == nil {
		fmt.Printf("No closest parcel found\n")
		return
	}
	simulation.board.At(uint(closestParcel.pos.X), uint(closestParcel.pos.Y)).Blocked = false
	path, err := pathfinding.Resolve(&simulation.board, f.pos, closestParcel.pos)
	simulation.board.At(uint(closestParcel.pos.X), uint(closestParcel.pos.Y)).Blocked = true
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}
	if len(path) <= 1 {
		// TODO: pick up parcel
		return
	}
	f.pos = path[0]
}
