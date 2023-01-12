package main

import "fmt"

type gridUnit uint32

type simulation struct {
	cycle     uint32
	warehouse warehouse
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

type entity interface {
	Name() stringer
	Kind() string
	Coord() coordinate
}

type warehouse struct {
	length gridUnit
	width  gridUnit

	parcels   []parcel
	forklifts []forklift
	trucks    []truck
}

type parcel struct {
	name string
	coordinate
	weight weight
}

func (parcel parcel) Name() stringer {
	return stringer(parcel.name)
}

func (parcel parcel) Kind() string {
	return "parcel"
}

func (parcel parcel) Coord() coordinate {
	return parcel.coordinate
}

type forklift struct {
	name string
	coordinate
}

func (forklift forklift) Name() stringer {
	return stringer(forklift.name)
}

func (forklift forklift) Kind() string {
	return "forklift"
}

func (forklift forklift) Coord() coordinate {
	return forklift.coordinate
}

type truck struct {
	name string
	coordinate
	maxWeight weight
	available uint32
}

func (truck truck) Name() stringer {
	return stringer(truck.name)
}

func (truck truck) Kind() string {
	return "truck"
}

func (truck truck) Coord() coordinate {
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
