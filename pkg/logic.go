package pkg

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

// Logic entrypoint for running the algorithm and update the UI
func Logic(gofeur *Gofeur) {
	ui := gofeur.ui
	str := "Go "
	count := 0

	for {
		time.Sleep(1 * time.Second)
		if count%2 == 0 {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s QUOI?", str))
		} else {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s FEUR...", str))
		}
		if !gofeur.ui.historic.IsRowSelected {
			if count == 2 {
				ui.DumpActionInStateBox(gofeur.sb.Forklifts[0], "TEST")
			} else {
				ui.DumpActionInStateBox(gofeur.sb.Trucks[0], "WAINTING")
			}
			ui.OutputBox.SetCell(count, 0, tview.NewTableCell(fmt.Sprintf("round %d\n", count)))
			count++
		}
		ui.App.Draw()
	}
}
