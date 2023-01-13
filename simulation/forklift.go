package simulation

import (
	"errors"

	"github.com/adrienlucbert/gofeur/logger"
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
	status ForkLiftStatus
	target optional.Optional[prop]
	path   optional.Optional[[]pkg.Vector]
}

func newForkliftFromParsing(from *parsing.Forklift) forklift {
	return forklift{
		name:   from.Name,
		pos:    pkg.Vector{X: int(from.X), Y: int(from.Y)},
		parcel: optional.NewEmpty[*parcel](),
		status: Empty,
		target: optional.NewEmpty[prop](),
		path:   optional.NewEmpty[[]pkg.Vector](),
	}
}

var errParcelNotFound = errors.New("No closest parcel found")
var errTruckNotFound = errors.New("No closest truck found")

func (f *forklift) findPathToTarget(simulation *Simulation) error {
	simulation.board.At(uint(f.target.Value().Pos().X), uint(f.target.Value().Pos().Y)).Blocked = false
	path, err := pathfinding.Resolve(&simulation.board, f.pos, f.target.Value().Pos())
	simulation.board.At(uint(f.target.Value().Pos().X), uint(f.target.Value().Pos().Y)).Blocked = true
	if err != nil {
		return err
	}
	f.path.Set(path)
	return nil
}

func (f *forklift) findClosestParcel(simulation *Simulation) error {
	// PERF: don't refetch target if not reached
	if target := findClosestParcel(simulation.parcels, f.pos); target != nil {
		f.target.Set(target)
	} else {
		f.target.Clear()
		return errParcelNotFound
	}
	return f.findPathToTarget(simulation)
}

func (f *forklift) findClosestTruck(simulation *Simulation) error {
	// PERF: don't refetch target if not reached
	if target := findClosestTruck(simulation.trucks, f.pos); target != nil {
		f.target.Set(target)
	} else {
		f.target.Clear()
		return errTruckNotFound
	}
	return f.findPathToTarget(simulation)
}

var errForkliftAlreadyLoaded = errors.New("Forklift already loaded")

func (f *forklift) grabParcel(parcel *parcel) error {
	if f.status == Loaded {
		return errForkliftAlreadyLoaded
	}
	f.target.Clear()
	f.path.Clear()
	f.parcel.Set(parcel)
	f.status = Loaded
	parcel.status = Carried
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
	f.parcel.Value().status = DroppedOff
	truck.load += f.parcel.Value().weight
	f.target.Clear()
	f.path.Clear()
	f.parcel.Clear()
	f.status = Empty
	return nil
}

func (f *forklift) seekParcel(simulation *Simulation) {
	if !f.target.HasValue() || simulation.board.At(uint(f.path.Value()[0].X), uint(f.path.Value()[0].Y)).Blocked {
		if err := f.findClosestParcel(simulation); err != nil {
			logger.Error("%s\n", err.Error())
			return
		}
	}
	if len(f.path.Value()) <= 1 {
		if err := f.grabParcel(f.target.Value().(*parcel)); err != nil {
			logger.Error("%s\n", err.Error())
		}
		return
	}
	f.pos = f.path.Value()[0]
	f.path.Set(f.path.Value()[1:])
}

func (f *forklift) seekTruck(simulation *Simulation) {
	if !f.target.HasValue() || simulation.board.At(uint(f.path.Value()[0].X), uint(f.path.Value()[0].Y)).Blocked {
		if err := f.findClosestTruck(simulation); err != nil {
			logger.Error("%s\n", err.Error())
			return
		}
	}
	if len(f.path.Value()) <= 1 {
		if err := f.depositParcel(f.target.Value().(*truck)); err != nil {
			logger.Error("%s\n", err.Error())
		}
		return
	}
	f.pos = f.path.Value()[0]
	f.path.Set(f.path.Value()[1:])
}

func (f *forklift) simulateRound(simulation *Simulation) {
	switch f.status {
	case Empty:
		f.seekParcel(simulation)
	case Loaded:
		f.seekTruck(simulation)
	}
}
