package parsing

import (
	"fmt"
)

type Color uint8

const (
	Yellow Color = iota
	Green
	Blue
)

type Startup struct {
	Width  uint
	Length uint
	Rounds uint
}

type StorageBuilding struct {
	Forklifts []Forklift
	Packs     []Parcel
	Trucks    []Truck
}

type Forklift struct {
	Name string
	X    uint
	Y    uint
}

type Truck struct {
	Name   string
	X      uint
	Y      uint
	Weight uint
	RAvail uint
}

func (t Truck) String() string {
	return "T"
}

type Parcel struct {
	Name  string
	X     uint
	Y     uint
	Color Color
}

func (t Forklift) String() string {
	return "F"
}

func (p Parcel) String() string {
	var str string
	white := "[white]"

	switch p.Color {
	case Green:
		str = "[green]"
	case Yellow:
		str = "[yellow]"
	case Blue:
		str = "[blue]"
	}
	return fmt.Sprintf("%sP%s", str, white)
}

type Gofeur struct {
	ST Startup
	SB StorageBuilding
}
