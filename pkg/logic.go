package pkg

import (
	"time"
	"fmt"
)

// Logic entrypoint for running the algorithm and update the UI
func Logic(ui *UI) {
	str := "Go "
	count := 0
	for {
		time.Sleep(4 * time.Second)
		if count%2 == 0 {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s QUOI?", str))
		} else {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s FEUR...", str))
		}
		fmt.Fprintf(ui.OutputBox, "round: %d\n", count)
		count++
		ui.App.Draw()
	}
}
