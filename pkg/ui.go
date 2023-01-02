package pkg

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	App                  *tview.Application
	StorageBuildingTable *tview.Table
	InfoBox              *tview.TextView
	OutputBox            *tview.TextView
	Layout               *tview.Flex
	building             [][]any
}

func UIStart() *UI {
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
		}
	}
	// This is only for testing purpose
	ui.building[sb.Packs[0].Y][sb.Packs[0].X] = sb.Packs[0]
	ui.building[sb.Packs[1].Y][sb.Packs[1].X] = sb.Packs[1]
	ui.building[sb.Transpals[0].Y][sb.Transpals[0].X] = sb.Transpals[0]
	ui.building[sb.Transpals[1].Y][sb.Transpals[1].X] = sb.Transpals[1]
	ui.building[sb.Trucks[0].Y][sb.Trucks[0].X] = sb.Trucks[0]

	ui.initUI()
	ui.Update(li, ly)
	return ui
}

func (ui *UI) initUI() {
	storageBuildingTable := tview.NewTable().
		SetSelectable(true, true).
		SetBorders(true)

	// This will display infos of each components presents in the storageTable, such as;
	// truck state (WAITING, GONE), transpals actions (GO, WAIT, TAKE, LEAVE) etc...
	infoBox := tview.NewTextView().
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			ui.App.Draw()
		})
	infoBox.
		SetBorder(true).
		SetTitle("Infos")

		// This will show the expected output
	outputBox := tview.NewTextView().
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			ui.App.Draw()
		})
	outputBox.
		SetBorder(true).
		SetTitle("Output")

	globalLayout := tview.NewFlex().
		AddItem(outputBox, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(storageBuildingTable, 0, 2, true).
			AddItem(infoBox, 5, 1, false), 0, 3, false)
	globalLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			ui.App.Stop()
			return nil
		}
		return event
	})

	ui.StorageBuildingTable = storageBuildingTable
	ui.InfoBox = infoBox
	ui.OutputBox = outputBox
	ui.Layout = globalLayout

	storageBuildingTable.SetSelectionChangedFunc(func(col, row int) {
        ui.InfoBox.Clear()
		fmt.Fprintf(ui.InfoBox, "x: %d, y: %d", col, row)
	})
}

// Don't mind this func it will be usefull
func (ui *UI) Update(x, y int) {
	color := tcell.ColorWhite

	for row := 0; row < y; row++ {
		for col := 0; col < x; col++ {
			ui.StorageBuildingTable.SetCell(row, col,
				tview.NewTableCell(fmt.Sprint(ui.building[col][row])).
					SetTextColor(color).
					SetExpansion(1).
					SetAlign(tview.AlignCenter))
		}
	}
	ui.Render()
}

func (ui *UI) Render() {
	ui.App.ForceDraw()
}
