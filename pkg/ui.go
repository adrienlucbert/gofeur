package pkg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type history struct {
	book          []string
	pastSelection int
	IsRowSelected bool
}

// UI est la struct regroupant tous les composants du TUI
type UI struct {
	App                  *tview.Application
	StorageBuildingTable *tview.Table
	InfoBox              *tview.TextView
	OutputBox            *tview.Table
	StateBox             *tview.TextView
	Layout               *tview.Flex
	building             [][]any
	historic             history
}

func addElementsToBuilding[T any](elements []T, building [][]any) {
	for _, elem := range elements {
		va := reflect.ValueOf(&elem).Elem()
		fieldX := va.FieldByName("X").Uint()
		fieldY := va.FieldByName("Y").Uint()
		building[fieldY][fieldX] = elem
	}
}

// UIStart instantiate UI Application and setup its environment
func UIStart(st Startup, sb StorageBuilding) *UI {
	app := tview.NewApplication()

	ui := &UI{
		App: app,
	}

	w := int(st.Width)
	l := int(st.Length)

	ui.building = make([][]any, l)

	for y := 0; y < l; y++ {
		for i := 0; i < w; i++ {
			ui.building[y] = append(ui.building[y], ".")
		}
	}
	addElementsToBuilding(sb.Packs, ui.building)
	addElementsToBuilding(sb.Forklifts, ui.building)
	addElementsToBuilding(sb.Trucks, ui.building)

	ui.initUI()
	ui.updateStorageBuildingTable(w, l)
	ui.updateOutputBox()
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

	stateBox := tview.NewTextView().
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			ui.App.Draw()
		})

	stateBox.
		SetBorder(true).
		SetTitle("State")
		// Add a mouse event handler to handle mouse events

	outputBox := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	outputBox.
		SetBorder(true).
		SetTitle("Output")

	globalLayout := tview.NewFlex().
		AddItem(outputBox, 15, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(storageBuildingTable, 0, 2, true).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(infoBox, 0, 1, false).
				AddItem(stateBox, 0, 1, false), 0, 1, false), 0, 5, false)
	globalLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			ui.App.Stop()
			return nil
		}
		if event.Key() == tcell.KeyRight {
			ui.historic.IsRowSelected = false
			ui.StateBox.Clear()
			pastHistory := strings.Join(ui.historic.book, "")
			fmt.Fprint(ui.StateBox, pastHistory)
		}
		return event
	})

	ui.StorageBuildingTable = storageBuildingTable
	ui.InfoBox = infoBox
	ui.OutputBox = outputBox
	ui.StateBox = stateBox
	ui.Layout = globalLayout
	ui.historic.IsRowSelected = false
}

// DumpActionInStateBox Display the given element and action in the state box
func (ui *UI) DumpActionInStateBox(elem any, action string) {
	fmt.Fprintf(ui.StateBox, "%s %s\n", fmt.Sprint(elem), action)
	actuState := fmt.Sprintf("%s %s\n", fmt.Sprint(elem), action)
	ui.historic.book = append(ui.historic.book, actuState)
}

func (ui *UI) updateOutputBox() {
	ui.OutputBox.SetSelectionChangedFunc(func(row, col int) {
		ui.StateBox.Clear()
		fmt.Fprintf(ui.StateBox, "row %s\n", ui.historic.book[row])
		ui.historic.IsRowSelected = true
		ui.historic.pastSelection = row
	})
}

func (ui *UI) updateStorageBuildingTable(x, y int) {
	color := tcell.ColorWhite

	ui.StorageBuildingTable.SetSelectionChangedFunc(func(col, row int) {
		if col >= x || row >= y || row < 0 || col < 0 {
			return
		}
		ui.InfoBox.Clear()
		if ui.building[col][row] == "." {
			fmt.Fprintf(ui.InfoBox, "x: %d, y: %d", col, row)
			return
		}
		va := reflect.ValueOf(ui.building[col][row])
		name := va.FieldByName("Name").String()
		fmt.Fprintf(ui.InfoBox, "Name: %s, x: %d, y: %d", name, col, row)
	})
	for row := 0; row < y; row++ {
		for col := 0; col < x; col++ {
			ui.StorageBuildingTable.SetCell(row, col,
				tview.NewTableCell(fmt.Sprint(ui.building[row][col])).
					SetTextColor(color).
					SetExpansion(1).
					SetAlign(tview.AlignCenter))
		}
	}
	ui.render()
}

func (ui *UI) render() {
	ui.App.ForceDraw()
}
