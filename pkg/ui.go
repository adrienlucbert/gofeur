package pkg

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	App                  *tview.Application
	StorageBuildingTable *tview.Table
	InfoBox              *tview.Box
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
	// This is temporary (Wainting for the parser to be done)
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
	storageBuildingTable := tview.NewTable().
		SetBorders(true)

		// This will display infos of each components presents in the storageTable, such as;
		// truck state (WAITING, GONE), transpals actions (GO, WAIT, TAKE, LEAVE) etc...
	infoBox := tview.NewBox().
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
}

// don't mind  this func it will be usefull
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
