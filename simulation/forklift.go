package simulation

import (
	"errors"
	"fmt"

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
	// Grabbing is the forklift's state when it's about to take a parcel
	Grabbing
	// Dropping is the forklift's state when it's about to drop a parcel
	Dropping
	// Loaded is the forklift's state when it carries a parcel
	Loaded
)

type forkliftAction interface {
	fmt.Stringer
}

type forkliftWaitAction struct{}

func (a forkliftWaitAction) String() string {
	return "WAIT"
}

type forkliftGoAction struct {
	pos pkg.Vector
}

func (a forkliftGoAction) String() string {
	return fmt.Sprintf("GO [%d,%d]", a.pos.X, a.pos.Y)
}

type forkliftTakeAction struct {
	parcel *parcel
}

func (a forkliftTakeAction) String() string {
	return fmt.Sprintf("TAKE %s %s", a.parcel.name, a.parcel.color)
}

type forkliftLeaveAction struct {
	parcel *parcel
}

func (a forkliftLeaveAction) String() string {
	return fmt.Sprintf("LEAVE %s %s", a.parcel.name, a.parcel.color)
}

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

var (
	errParcelNotFound = errors.New("No closest parcel found")
	errTruckNotFound  = errors.New("No closest truck found")
)

type pathToTargetError struct {
	pathfinding error
}

func (err pathToTargetError) Error() string {
	return fmt.Sprintf("No path to target: %s", err.pathfinding.Error())
}

func (f *forklift) findPathToTarget(simulation *Simulation) error {
	simulation.board.At(uint(f.target.Value().Pos().X), uint(f.target.Value().Pos().Y)).Blocked = false
	path, err := pathfinding.Resolve(&simulation.board, f.pos, f.target.Value().Pos())
	simulation.board.At(uint(f.target.Value().Pos().X), uint(f.target.Value().Pos().Y)).Blocked = true
	if err != nil {
		return pathToTargetError{pathfinding: err}
	}
	f.path.Set(path)
	return nil
}

func (f *forklift) findClosestParcel(simulation *Simulation) error {
	// PERF: don't refetch target if not reached
	if target := findClosestParcel(simulation.parcels, f.pos); target != nil {
		target.status = Targeted
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

func (f *forklift) startGrabbingParcel() error {
	if f.status == Loaded {
		return errForkliftAlreadyLoaded
	}
	f.status = Grabbing
	return nil
}

func (f *forklift) finishGrabbingParcel() {
	f.parcel.Set(f.target.Value().(*parcel))
	f.parcel.Value().status = Carried
	f.target.Clear()
	f.path.Clear()
	f.status = Loaded
}

func (f *forklift) dropParcelFocus() {
	f.target.Value().(*parcel).status = StandingBy
	f.target.Clear()
	f.path.Clear()
}

var (
	errForkliftEmpty = errors.New("Forklift is empty")
	errTruckFull     = errors.New("Truck is full")
)

func (f *forklift) startDroppingParcel() error {
	if !f.parcel.HasValue() {
		return errForkliftEmpty
	}
	truck := f.target.Value().(*truck)
	if truck.load+f.parcel.Value().weight > truck.capacity {
		return errTruckFull
	}
	f.status = Dropping
	return nil
}

func (f *forklift) finishDroppingParcel() {
	truck := f.target.Value().(*truck)
	truck.load += f.parcel.Value().weight
	f.target.Clear()
	f.path.Clear()
	f.parcel.Value().status = DroppedOff
	f.parcel.Clear()
	f.status = Empty
}

func (f *forklift) seekParcel(simulation *Simulation) forkliftAction {
	if !f.target.HasValue() || simulation.board.At(uint(f.path.Value()[0].X), uint(f.path.Value()[0].Y)).Blocked {
		if f.target.HasValue() {
			f.dropParcelFocus()
		}
		if err := f.findClosestParcel(simulation); err != nil {
			logger.Debug("%s\n", err.Error())
			return forkliftWaitAction{}
		}
	}
	if len(f.path.Value()) <= 1 {
		if err := f.startGrabbingParcel(); err != nil {
			logger.Debug("%s\n", err.Error())
		}
		return forkliftTakeAction{f.target.Value().(*parcel)}
	}
	dest := f.path.Value()[0]
	f.pos = dest
	f.path.Set(f.path.Value()[1:])
	return forkliftGoAction{dest}
}

func (f *forklift) seekTruck(simulation *Simulation) forkliftAction {
	if !f.target.HasValue() || !f.target.Value().IsAvailable() || simulation.board.At(uint(f.path.Value()[0].X), uint(f.path.Value()[0].Y)).Blocked {
		if err := f.findClosestTruck(simulation); err != nil {
			logger.Debug("%s\n", err.Error())
			return forkliftWaitAction{}
		}
	}
	if len(f.path.Value()) <= 1 {
		if err := f.startDroppingParcel(); err != nil {
			logger.Debug("%s\n", err.Error())
		}
		return forkliftLeaveAction{f.parcel.Value()}
	}
	dest := f.path.Value()[0]
	f.pos = dest
	f.path.Set(f.path.Value()[1:])
	return forkliftGoAction{dest}
}

func (f *forklift) simulateRound(simulation *Simulation) {
	var action forkliftAction
	switch f.status {
	case Grabbing:
		f.finishGrabbingParcel()
	case Dropping:
		f.finishDroppingParcel()
	}
	switch f.status {
	case Empty:
		action = f.seekParcel(simulation)
	case Grabbing:
	case Loaded:
		action = f.seekTruck(simulation)
	}
	logger.Info("%s %s\n", f.name, action.String())
}
