// Package pkg implements everything relative to the project and so on and so forthgofgo
package pkg

import (
	"fmt"
)

type color uint8

// color iota
const (
	Green = iota
	Yellow
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
	Color color
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
	ui *UI
	st Startup
	sb StorageBuilding
}

// Init the program and UI
func (gofeur *Gofeur) Init() {
	gofeur.ui = UIStart(gofeur.st, gofeur.sb)
}

// RunUI the UI
func (gofeur *Gofeur) RunUI() {
	feurUI := gofeur.ui

	if err := feurUI.App.SetRoot(feurUI.Layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
