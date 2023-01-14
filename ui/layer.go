package ui

import (
	"fmt"
	"time"

	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/simulation"
	"github.com/rivo/tview"
)

// Layer is an optional application layer responsible for displaying the UI
type Layer struct {
	Gofeur     *parsing.Simulation
	Simulation *simulation.Simulation
	ui         *UI
}

func (layer *Layer) run() {
	err := layer.ui.App.SetRoot(layer.ui.Layout, true).
		EnableMouse(true).
		Run()
	if err != nil {
		panic(err)
	}
}

// Attach initializes the UILayer
func (layer *Layer) Attach() {
	layer.ui = Start(layer.Gofeur)
	go layer.run()
}

// Update updates the UI and re-renders it
func (layer *Layer) Update(elapsedTime time.Duration) {
	if layer.Simulation.Round%2 == 0 {
		layer.ui.OutputBox.SetTitle("Go QUOI?")
	} else {
		layer.ui.OutputBox.SetTitle("Go FEUR...")
	}
	if !layer.ui.historic.IsRowSelected {
		if layer.Simulation.Round == 2 {
			layer.ui.DumpActionInStateBox(layer.Gofeur.Warehouse.Forklifts[0], "TEST")
		} else {
			layer.ui.DumpActionInStateBox(layer.Gofeur.Warehouse.Trucks[0], "WAITING")
		}
		layer.ui.OutputBox.SetCell(int(layer.Simulation.Round), 0, tview.NewTableCell(fmt.Sprintf("round %d\n", layer.Simulation.Round)))
	}
	layer.ui.App.Draw()
}

// Detach dismounts the UILayer
func (layer *Layer) Detach() {
	layer.ui.App.Stop()
}
