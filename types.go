package main

type Unit uint32

type Simulation struct {
	cycle     uint32
	warehouse Warehouse
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

type Forklift struct {
	name string
	Coordinate
}

type Truck struct {
	name string
	Coordinate
	max_weight Weight
	available  uint32
}

type Weight uint32

type Coordinate struct {
	X Unit
	Y Unit
}

const (
	yellow Weight = 100
	green  Weight = 200
	blue   Weight = 500
)
