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

// Startup struct define the 1st line of the file given to the program
type Startup struct {
	Width  uint
	Length uint
	Rounds uint
}

// StorageBuilding struct is a basic representation of the environment
type StorageBuilding struct {
	Forklifts []Forklift
	Packs     []Parcel
	Trucks    []Truck
}

// Forklift struct
type Forklift struct {
	Name string
	X    uint
	Y    uint
}

// Truck struct
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

// Parcel struct
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

// Gofeur base struct that is logic entrypoint to the other struct of the program
type Gofeur struct {
	ST Startup
	SB StorageBuilding
}
