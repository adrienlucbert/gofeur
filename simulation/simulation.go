// Package simulation holds the logic behind a warehouse simulation
package simulation

import (
	"github.com/adrienlucbert/gofeur/board"
	"github.com/adrienlucbert/gofeur/logger"
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
)

// Status represents the simulation's possible states
type Status int

const (
	// Idle is Gofeur's state before it started
	Idle Status = iota
	// Running is Gofeur's state before it's running
	Running
	// Finished is Gofeur's state when it's over and no parcel remains in the warehouse
	Finished
	// Unfinished is Gofeur's state when it's over but parcels remains in the warehouse
	Unfinished
)

type prop interface {
	Pos() pkg.Vector
}

// Simulation represents the simulation data
type Simulation struct {
	MaxRound  uint
	Round     uint
	Status    Status
	board     board.Board
	forklifts []forklift
	parcels   []parcel
	trucks    []truck
}

// IsRunning returns whether or not the simulation is in the Running state
func (s *Simulation) IsRunning() bool {
	return s.Status == Running
}

func findClosestParcel(parcels []parcel, pos pkg.Vector) *parcel {
	var closestParcel *parcel
	var closestParcelDistance float32
	for i := range parcels {
		if parcels[i].status != StandingBy {
			continue
		}
		parcel := &parcels[i]
		parcelDistance := pos.SquaredDistance(parcel.pos)
		if closestParcel == nil || parcelDistance < closestParcelDistance {
			closestParcel = parcel
			closestParcelDistance = parcelDistance
		}
	}
	return closestParcel
}

// TODO: make truck and parcel implement a common interface to avoid repeating code
// NOTE: turns out the type-specific filter condition makes it difficult as
// predicates can't be called with an interface as parameter
func findClosestTruck(trucks []truck, pos pkg.Vector) *truck {
	var closestTruck *truck
	var closestTruckDistance float32
	for i := range trucks {
		// TODO: stop if truck is gone
		truck := &trucks[i]
		truckDistance := pos.SquaredDistance(truck.pos)
		if closestTruck == nil || truckDistance < closestTruckDistance {
			closestTruck = truck
			closestTruckDistance = truckDistance
		}
	}
	return closestTruck
}

func (s *Simulation) start() {
	s.Round = 0
	s.Status = Running
}

// New initializes a Simulation object
func New(gofeur *parsing.Gofeur) Simulation {
	s := Simulation{}
	s.MaxRound = gofeur.ST.Rounds
	s.board = board.New(gofeur.ST.Width, gofeur.ST.Length)
	for i := range gofeur.SB.Forklifts {
		s.forklifts = append(s.forklifts, newForkliftFromParsing(&gofeur.SB.Forklifts[i]))
	}
	for i := range gofeur.SB.Packs {
		s.parcels = append(s.parcels, newParcelFromParsing(&gofeur.SB.Packs[i]))
	}
	for i := range gofeur.SB.Trucks {
		s.trucks = append(s.trucks, newTruckFromParsing(&gofeur.SB.Trucks[i]))
	}
	s.updateBoard()
	return s
}

func (s *Simulation) updateBoard() {
	s.board.Clear()
	for i := range s.parcels {
		if s.parcels[i].status != StandingBy {
			continue
		}
		s.board.At(uint(s.parcels[i].pos.X), uint(s.parcels[i].pos.Y)).Blocked = true
		s.board.At(uint(s.parcels[i].pos.X), uint(s.parcels[i].pos.Y)).DebugChar = map[uint]rune{
			100: '1',
			200: '2',
			500: '3',
		}[s.parcels[i].weight]
	}
	for i := range s.forklifts {
		s.board.At(uint(s.forklifts[i].pos.X), uint(s.forklifts[i].pos.Y)).Blocked = true
		s.board.At(uint(s.forklifts[i].pos.X), uint(s.forklifts[i].pos.Y)).DebugChar = 'L'
	}
	for i := range s.trucks {
		s.board.At(uint(s.trucks[i].pos.X), uint(s.trucks[i].pos.Y)).Blocked = true
		s.board.At(uint(s.trucks[i].pos.X), uint(s.trucks[i].pos.Y)).DebugChar = 'T'
	}
}

func (s *Simulation) areAnyParcelsLeft() bool {
	for i := range s.parcels {
		if s.parcels[i].status != DroppedOff {
			return true
		}
	}
	return false
}

func (s *Simulation) simulateRound() {
	if !s.areAnyParcelsLeft() {
		s.Status = Finished
		return
	}
	logger.Debug("Round %d\n", s.Round+1)
	for i := range s.forklifts {
		s.forklifts[i].simulateRound(s)
	}
	for i := range s.trucks {
		s.trucks[i].simulateRound(s)
	}
	s.updateBoard()
	logger.Debug("%s\n", s.board.String())

	// Increment round and end simulation if needed
	s.Round++
	if s.Round >= s.MaxRound {
		s.Status = Unfinished
	}
}

func (s *Simulation) terminate() {
}
