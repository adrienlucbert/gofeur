package parsing

import "fmt"

type gridUnit uint32

// SimulationCycle represents the number of cycle the simulation should be run
type SimulationCycle uint32

// Simulation contains the content of a parsed input file.
type Simulation struct {
	Cycle     SimulationCycle
	Warehouse Warehouse
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

type entity interface {
	stringerName() stringer
	kind() string
	coord() coordinate
}

// Warehouse represents stores parsed properties and entities
type Warehouse struct {
	Length gridUnit
	Width  gridUnit

	Parcels   []Parcel
	Forklifts []Forklift
	Trucks    []Truck
}

// Parcel represents a parsed parcel
type Parcel struct {
	Name string
	coordinate
	Weight weight
}

func (parcel Parcel) stringerName() stringer {
	return stringer(parcel.Name)
}

func (parcel Parcel) kind() string {
	return "parcel"
}

func (parcel Parcel) coord() coordinate {
	return parcel.coordinate
}

// Forklift represents a parsed forklift
type Forklift struct {
	Name string
	coordinate
}

func (forklift Forklift) stringerName() stringer {
	return stringer(forklift.Name)
}

func (forklift Forklift) kind() string {
	return "forklift"
}

func (forklift Forklift) coord() coordinate {
	return forklift.coordinate
}

// Truck represents a parsed truck
type Truck struct {
	Name string
	coordinate
	MaxWeight weight
	Available uint32
}

func (truck Truck) stringerName() stringer {
	return stringer(truck.Name)
}

func (truck Truck) kind() string {
	return "truck"
}

func (truck Truck) coord() coordinate {
	return truck.coordinate
}

type weight uint32

type coordinate struct {
	X gridUnit
	Y gridUnit
}

func (coord coordinate) String() string {
	return fmt.Sprintf("(x: %d, y: %d)", coord.X, coord.Y)
}

const (
	yellow weight = 100
	green  weight = 200
	blue   weight = 500
)
