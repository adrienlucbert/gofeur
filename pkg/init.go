package pkg

import (
	"fmt"
)

type COLOR uint8

const (
	GREEN = iota
	BLUE
	RED
)

type Position struct {
	x int
	y int
}

type StorageBuilding struct {
	transpals []Transpalette
	packs     []Parcel
	truck     []Truck
}

type Transpalette struct {
	x int
	y int
}
type Truck struct {
	x int
	y int
}

func (t Truck) String() string {
	return fmt.Sprintf("Q")
}

type Parcel struct {
	x     int
	y     int
	color COLOR
}

func (t Transpalette) String() string {
	return fmt.Sprintf("T")
}

func (p Parcel) String() string {
	var str string
	white := "[white]"

	switch p.color {
	case GREEN:
		str = "[green]"
	case RED:
		str = "[red]"
	case BLUE:
		str = "[blue]"
	}
	return fmt.Sprintf("%sP%s", str, white)
}

// Temporary (Wainting for the parser to be done)
func InitStorageBuilding() *StorageBuilding {
	t1 := Truck{3, 10}
	truck := []Truck{t1}
	trans := []Transpalette{{0, 0}, {5, 2}}
	parcels := []Parcel{
		{2, 2, GREEN},
		{3, 5, RED},
	}
	sb := StorageBuilding{trans, parcels, truck}
	return &sb
}
