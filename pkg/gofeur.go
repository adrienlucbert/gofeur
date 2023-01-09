package pkg

import (
	"fmt"
)

type Color uint8

const (
	Green = iota
	Yellow
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
	return fmt.Sprintf("T")
}

type Parcel struct {
	Name  string
	X     uint
	Y     uint
	Color Color
}

func (t Forklift) String() string {
	return fmt.Sprintf("F")
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
	Ui *UI
	st Startup
	sb StorageBuilding
}

func (gofeur *Gofeur) Init() {
	gofeur.Ui = UIStart(gofeur.st, gofeur.sb)
}

func (gofeur *Gofeur) Run() {
	feurUI := gofeur.Ui

	if err := feurUI.App.SetRoot(feurUI.Layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
