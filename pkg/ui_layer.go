package pkg

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

// UILayer is an optional application layer responsible for displaying the UI
type UILayer struct {
	Gofeur *Gofeur
	ui     *UI
}

func (layer *UILayer) run() {
	err := layer.ui.App.SetRoot(layer.ui.Layout, true).
		EnableMouse(true).
		Run()
	if err != nil {
		panic(err)
	}
}

// Attach initializes the UILayer
func (layer *UILayer) Attach() {
	layer.ui = UIStart(layer.Gofeur.st, layer.Gofeur.sb)
	go layer.run()
}

// Update updates the UI and re-renders it
func (layer *UILayer) Update(elapsedTime time.Duration) {
	if layer.Gofeur.Step%2 == 0 {
		layer.ui.OutputBox.SetTitle("Go QUOI?")
	} else {
		layer.ui.OutputBox.SetTitle("Go FEUR...")
	}
	layer.ui.OutputBox.SetCell(int(layer.Gofeur.Step), 0, tview.NewTableCell(fmt.Sprintf("round %d\n", layer.Gofeur.Step)))
	layer.ui.App.Draw()
}

// Detach dismounts the UILayer
func (layer *UILayer) Detach() {
	layer.ui.App.Stop()
}
