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
	Coordonate
	weight Weight
}

type Forklift struct {
	name string
	Coordonate
}

type Truck struct {
	name string
	Coordonate
	max_weight Weight
	available  uint32
}

type Weight uint32

type Coordonate struct {
	X Unit
	Y Unit
}

const (
	yellow Weight = 100
	green  Weight = 200
	blue   Weight = 500
)
