package pkg

import (
	"fmt"

	"github.com/adrienlucbert/gofeur/pathfinding/board"
)

// SimulationStatus represents the simulation's possible states
type SimulationStatus int

const (
	// Idle is Gofeur's state before it started
	Idle SimulationStatus = iota
	// Running is Gofeur's state before it's running
	Running
	// Finished is Gofeur's state when it's over and no parcel remains in the warehouse
	Finished
	// Unfinished is Gofeur's state when it's over but parcels remains in the warehouse
	Unfinished
)

// Simulation represents the simulation data
type Simulation struct {
	MaxRound  uint
	Round     uint
	Status    SimulationStatus
	board     board.Board
	forklifts []Forklift
	parcels   []Parcel
	trucks    []Truck
}

func (s *Simulation) start() {
	s.Round = 0
	s.Status = Running
}

// NewSimulation initializes a Simulation object
func NewSimulation(gofeur *Gofeur) Simulation {
	s := Simulation{}
	s.MaxRound = gofeur.st.Rounds
	s.board = board.New(gofeur.st.Width, gofeur.st.Length)
	s.forklifts = gofeur.sb.Forklifts
	s.parcels = gofeur.sb.Packs
	s.trucks = gofeur.sb.Trucks
	s.updateBoard()
	return s
}

func (s *Simulation) updateBoard() {
	s.board.Clear()
	for i := range s.parcels {
		s.board.At(s.parcels[i].X, s.parcels[i].Y).Blocked = true
		s.board.At(s.parcels[i].X, s.parcels[i].Y).DebugChar = map[color_t]rune{
			Green:  '1',
			Yellow: '2',
			Blue:   '3',
		}[s.parcels[i].Color]
	}
	for i := range s.forklifts {
		s.board.At(s.forklifts[i].X, s.forklifts[i].Y).Blocked = true
		s.board.At(s.forklifts[i].X, s.forklifts[i].Y).DebugChar = 'L'
	}
	for i := range s.trucks {
		s.board.At(s.trucks[i].X, s.trucks[i].Y).Blocked = true
		s.board.At(s.trucks[i].X, s.trucks[i].Y).DebugChar = 'T'
	}
}

func (s *Simulation) simulateRound() {
	fmt.Printf("Round %d\n", s.Round+1)
	for i := range s.forklifts {
		s.simulateForkliftBehaviour(&s.forklifts[i])
	}
	for i := range s.trucks {
		s.simulateTruckBehaviour(&s.trucks[i])
	}
	s.updateBoard()
	fmt.Printf("%s\n", s.board.String())

	// Increment round and end simulation if needed
	s.Round++
	if s.Round >= s.MaxRound {
		s.Status = Unfinished
	}
}

func (s *Simulation) simulateForkliftBehaviour(forklift *Forklift) {
}

func (s *Simulation) simulateTruckBehaviour(truck *Truck) {
}

func (s *Simulation) terminate() {
}
