package pkg

import (
	"fmt"
	"reflect"
)

type COLOR uint8

const (
	GREEN = iota
	YELLOW
	BLUE
)

type Startup struct {
	Width  uint
	Length uint
	Rounds uint
}

type StorageBuilding struct {
	Transpals []Transpals
	Packs     []Parcel
	Trucks    []Truck
}

type Transpals struct {
	Name string
	X    uint
	Y    uint
}

type Truck struct {
	Name    string
	X       uint
	Y       uint
	Weight  uint
	R_avail uint
}

func (t Truck) String() string {
	return fmt.Sprintf("Q")
}

type Parcel struct {
	Name  string
	X     uint
	Y     uint
	Color COLOR
}

func (t Transpals) String() string {
	return fmt.Sprintf("T")
}

func (p Parcel) String() string {
	var str string
	white := "[white]"

	switch p.Color {
	case GREEN:
		str = "[green]"
	case YELLOW:
		str = "[yellow]"
	case BLUE:
		str = "[blue]"
	}
	return fmt.Sprintf("%sP%s", str, white)
}

type Gofeur struct {
	Ui *UI
	st Startup
	sb StorageBuilding
}

func (gofeur *Gofeur) setStartup(w, l, r uint) {
	val := reflect.ValueOf(gofeur)

	if val.Field(2).IsValid() == false {
		fmt.Println(`The startup has already been set
        , so we will keep the previous value and this line is not evaluated`)
		return
	}
}

func (gofeur *Gofeur) Init() {
	gofeur.Ui = UIStart()
}

func (gofeur *Gofeur) Run() {
	feurUI := gofeur.Ui

	if err := feurUI.App.SetRoot(feurUI.Layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}