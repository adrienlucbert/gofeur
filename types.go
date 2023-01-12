package main

import "fmt"

type Unit uint32

type Simulation struct {
	cycle     uint32
	warehouse Warehouse
}

type stringer string

func (s stringer) String() string {
	return (string)(s)
}

type Entity interface {
	Name() stringer
	Kind() string
	Coord() Coordinate
}

type Warehouse struct {
	length Unit
	width  Unit

	parcels   []Parcel
	forklifts []Forklift
	trucks    []Truck
}

type Parcel struct {
	name string
	Coordinate
	weight Weight
}

func (parcel Parcel) Name() stringer {
	return (stringer)(parcel.name)
}

func (parcel Parcel) Kind() string {
	return "parcel"
}

func (parcel Parcel) Coord() Coordinate {
	return parcel.Coordinate
}

type Forklift struct {
	name string
	Coordinate
}

func (forklift Forklift) Name() stringer {
	return (stringer)(forklift.name)
}

func (forklift Forklift) Kind() string {
	return "forklift"
}

func (forklift Forklift) Coord() Coordinate {
	return forklift.Coordinate
}

type Truck struct {
	name string
	Coordinate
	max_weight Weight
	available  uint32
}

func (truck Truck) Name() stringer {
	return (stringer)(truck.name)
}

func (truck Truck) Kind() string {
	return "truck"
}

func (truck Truck) Coord() Coordinate {
	return truck.Coordinate
}

type Weight uint32

type Coordinate struct {
	X Unit
	Y Unit
}

func (coord Coordinate) String() string {
	return fmt.Sprintf("(x: %d, y: %d)", coord.X, coord.Y)
}

const (
	yellow Weight = 100
	green  Weight = 200
	blue   Weight = 500
)
