package pkg

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

// Logic entrypoint for running the algorithm and update the UI
func Logic(ui *UI) {
	str := "Go "
	count := 0
	for {
		time.Sleep(1 * time.Second)
		if count%2 == 0 {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s QUOI?", str))
		} else {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s FEUR...", str))
		}
		ui.OutputBox.SetCell(count, 0, tview.NewTableCell(fmt.Sprintf("round %d\n", count)))
		count++
		ui.App.Draw()
	}
}
