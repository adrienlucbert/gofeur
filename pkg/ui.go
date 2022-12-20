package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	App                  *tview.Application
	StorageBuildingTable *tview.Table
	InfoBox              *tview.Box
	OutputBox            *tview.Box
	FlexLayout           *tview.Flex
	building             [][]any
}

type Gofeur struct {
	ui *UI
}

func Start() *UI {
	app := tview.NewApplication()

	ui := &UI{
		App: app,
	}
	li := 16
	ly := 16

	ui.building = make([][]any, ly)

	sb := InitStorageBuilding()

	for y := 0; y < ly; y++ {
		for i := 0; i < li; i++ {
			ui.building[y] = append(ui.building[y], ".")
			fmt.Print(ui.building[y][i])

		}
		fmt.Print("\n")
	}
	ui.building[sb.packs[0].y][sb.packs[0].x] = sb.packs[0]
	ui.building[sb.packs[1].y][sb.packs[1].x] = sb.packs[1]
	ui.building[sb.transpals[0].y][sb.transpals[0].x] = sb.transpals[0]
	ui.building[sb.transpals[1].y][sb.transpals[1].x] = sb.transpals[1]
	ui.building[sb.truck[0].y][sb.truck[0].x] = sb.truck[0]

	ui.initUI(li, ly)
	ui.Update(li, ly)
	return ui
}

func (ui *UI) initUI(x, y int) {
	storageBuildingUI := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(false).
		SetChangedFunc(func() {
			ui.App.Draw()
		})
	storageTable := tview.NewTable().
		SetBorders(true)
        infoBox := tview.NewBox().
		SetBorder(true).
		SetTitle("Bottom (5 rows)")
	outputBox := tview.NewBox().
		SetBorder(true).
		SetTitle("Round")

	color := tcell.ColorWhite
	for row := 0; row < y; row++ {
		for col := 0; col < x; col++ {
			storageTable.SetCell(row, col,
				tview.NewTableCell(".").
					SetTextColor(color).
					SetExpansion(1).
					SetAlign(tview.AlignCenter))
		}
	}

	storageBuildingUI.
		SetBorder(true).
		SetTitle("Storage Building")

	flexLayout := tview.NewFlex().
		AddItem(outputBox, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(storageTable, 0, 2, true).
			AddItem(infoBox, 5, 1, false), 0, 3, false)

	ui.StorageBuildingTable = storageTable
	ui.InfoBox = infoBox
	ui.OutputBox = outputBox
	ui.FlexLayout = flexLayout
}

func (ui *UI) Update(x, y int) {
	for row := 0; row < y; row++ {
		for col := 0; col < x; col++ {
			ui.StorageBuildingTable.SetCell(row, col,
				tview.NewTableCell(fmt.Sprint(ui.building[col][row])).
					SetExpansion(1).
					SetAlign(tview.AlignCenter))
		}
	}
	ui.Render()
}

func (ui *UI) Render() {
	ui.App.ForceDraw()
}
