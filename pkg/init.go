package ui

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

// func (p Parcel) String() string {
//     var str string
//     white := "\033[0;37m"
//
//     switch p.color {
//     case GREEN:
//         str = "\033[0;32m"
//     case RED:
//         str = "\033[0;31m"
//     case BLUE:
//         str = "\033[0;31m"
//     }
//     return fmt.Sprintf("%sP%s", str, white)
// }

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

// func dump(building [11][11]any) {
// 	li := 10
// 	ly := 10
//
// 	line := strings.Repeat("-", li*3)
// 	fmt.Println(line)
// 	for y := 0; y < ly; y++ {
// 		for i := 0; i < li; i++ {
// 			fmt.Printf("[%s]", building[y][i])
// 		}
// 		fmt.Print("\n")
// 	}
// 	fmt.Println(line)
// }

func display() {
	li := 10
	ly := 10

	var building [11][11]any

	sb := InitStorageBuilding()

	for y := 0; y < ly; y++ {
		for i := 0; i < li; i++ {
			building[y][i] = "."
		}
	}
	building[sb.packs[0].y][sb.packs[0].x] = sb.packs[0]
	building[sb.packs[1].y][sb.packs[1].x] = sb.packs[1]
	// dump(building)
}

// func InitUi(x, y int) {
// 	li := x
// 	ly := y
//
// 	var building [16][16]any
//
// 	ui := tview.NewApplication()
//
// 	storageBuildingUI := tview.NewTextView().
// 		SetDynamicColors(true).
// 		SetRegions(false).
// 		SetChangedFunc(func() {
// 			ui.Draw()
// 		})
//
// 	storageTable := tview.NewTable().
// 		SetBorders(true)
//
// 	infoBox := tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)")
//
// 	color := tcell.ColorWhite
//
// 	for row := 0; row < y; row++ {
// 		for col := 0; col < x; col++ {
// 			storageTable.SetCell(row, col,
// 				tview.NewTableCell(fmt.Sprintf("%s", building[col][row])).
// 					SetTextColor(color).
// 					SetExpansion(1).
// 					SetAlign(tview.AlignCenter))
// 		}
// 	}
//
// 	storageBuildingUI.
// 		SetBorder(true).
// 		SetTitle("Storage Building")
//
// 	flexLayout := tview.NewFlex().
// 		AddItem(tview.NewBox().SetBorder(true).SetTitle("Round"), 0, 1, false).
// 		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
// 			AddItem(storageTable, 0, 2, true).
// 			AddItem(infoBox, 5, 1, false), 0, 3, false)
//
// 	if err := ui.SetRoot(flexLayout, true).SetFocus(flexLayout).Run(); err != nil {
// 		panic(err)
// 	}
// }
