package simulation

import (
	"errors"
	"fmt"

	"github.com/adrienlucbert/gofeur/optional"
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pathfinding"
	"github.com/adrienlucbert/gofeur/pkg"
)

type forklift struct {
	name string
	pos  pkg.Vector
	// TODO: load uint
	// TODO: state [LOADED|EMPTY]
	target optional.Optional[*parcel]
	path   optional.Optional[[]pkg.Vector]
}

func newForkliftFromParsing(from *parsing.Forklift) forklift {
	return forklift{
		name: from.Name,
		pos:  pkg.Vector{X: int(from.X), Y: int(from.Y)},
	}
}

// errPathNotFound is returned when a path couldn't be found
var errParcelNotFound = errors.New("No closest parcel found")

func (f *forklift) findTarget(simulation *Simulation) error {
	f.target.Set(findClosestParcel(simulation.parcels, f.pos))
	if !f.target.HasValue() {
		return errParcelNotFound
	}
	simulation.board.At(uint(f.target.Value().pos.X), uint(f.target.Value().pos.Y)).Blocked = false
	path, err := pathfinding.Resolve(&simulation.board, f.pos, f.target.Value().pos)
	simulation.board.At(uint(f.target.Value().pos.X), uint(f.target.Value().pos.Y)).Blocked = true
	if err != nil {
		return err
	}
	f.path.Set(path)
	return nil
}

func (f *forklift) simulateRound(simulation *Simulation) {
	if !f.target.HasValue() {
		if err := f.findTarget(simulation); err != nil {
			fmt.Print(err.Error())
			return
		}
	}
	if len(f.path.Value()) <= 1 {
		// TODO: pick up parcel
		return
	}
	f.pos = f.path.Value()[0]
}
