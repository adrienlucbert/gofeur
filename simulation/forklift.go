package simulation

import (
	"errors"
	"fmt"

	"github.com/adrienlucbert/gofeur/optional"
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pathfinding"
	"github.com/adrienlucbert/gofeur/pkg"
)

// ForkLiftStatus represents the forklift's possible states
type ForkLiftStatus int

const (
	// Empty is the forklift's state when it carries no parcel
	Empty ForkLiftStatus = iota
	// Loaded is the forklift's state when it carries a parcel
	Loaded
)

type forklift struct {
	name   string
	pos    pkg.Vector
	parcel optional.Optional[*parcel]
	state  ForkLiftStatus
	target optional.Optional[*parcel]
	path   optional.Optional[[]pkg.Vector]
}

func newForkliftFromParsing(from *parsing.Forklift) forklift {
	return forklift{
		name:   from.Name,
		pos:    pkg.Vector{X: int(from.X), Y: int(from.Y)},
		parcel: optional.NewEmpty[*parcel](),
		state:  Empty,
		target: optional.NewEmpty[*parcel](),
		path:   optional.NewEmpty[[]pkg.Vector](),
	}
}

var errParcelNotFound = errors.New("No closest parcel found")

func (f *forklift) findTarget(simulation *Simulation) error {
	// PERF: don't refetch target if not reached
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

var errForkliftAlreadyLoaded = errors.New("Forklift already loaded")

func (f *forklift) grabParcel(parcel *parcel) error {
	if f.state == Loaded {
		return errForkliftAlreadyLoaded
	}
	f.target.Clear()
	f.path.Clear()
	f.parcel.Set(parcel)
	f.state = Loaded
	parcel.carried = true
	return nil
}

var errForkliftEmpty = errors.New("Forklift is empty")
var errTruckFull = errors.New("Truck is full")

func (f *forklift) depositParcel(truck *truck) error {
	if !f.parcel.HasValue() {
		return errForkliftEmpty
	}
	if truck.load+f.parcel.Value().weight > truck.capacity {
		return errTruckFull
	}
	truck.load += f.parcel.Value().weight
	f.parcel.Clear()
	return nil
}

func (f *forklift) simulateRound(simulation *Simulation) {
	switch f.state {
	case Empty:
		if !f.target.HasValue() || simulation.board.At(uint(f.path.Value()[0].X), uint(f.path.Value()[0].Y)).Blocked {
			if err := f.findTarget(simulation); err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}
		}
		if len(f.path.Value()) <= 1 {
			fmt.Print("GRAB\n")
			if err := f.grabParcel(f.target.Value()); err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}
			return
		}
		f.pos = f.path.Value()[0]
		f.path.Set(f.path.Value()[1:])
	case Loaded:
		// TODO: go to nearest available truck
	}
}
