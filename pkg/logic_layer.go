package pkg

import (
	"time"
)

// LogicLayer is the application layer responsible for managing the game logic
type LogicLayer struct {
	Gofeur *Gofeur
}

// Attach initializes the LogicLayer
func (layer *LogicLayer) Attach() {
	layer.Gofeur.Status = Running
}

// Update runs the game logic
func (layer *LogicLayer) Update(elapsedTime time.Duration) {
	time.Sleep(1 * time.Second)
	layer.Gofeur.Step++
	if layer.Gofeur.Step > 3 {
		layer.Gofeur.Status = Finished
	}
}

// Detach handles the game end
func (layer *LogicLayer) Detach() {
	// TODO: print end-of-game data
}
