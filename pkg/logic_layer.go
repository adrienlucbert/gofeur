package pkg

import (
	"time"
)

// LogicLayer is the application layer responsible for managing the game logic
type LogicLayer struct {
	Simulation *Simulation
}

// Attach initializes the LogicLayer
func (layer *LogicLayer) Attach() {
	layer.Simulation.start()
}

// Update runs the game logic
func (layer *LogicLayer) Update(elapsedTime time.Duration) {
	layer.Simulation.simulateRound()
}

// Detach handles the game end
func (layer *LogicLayer) Detach() {
	layer.Simulation.terminate()
	// TODO: print end-of-game data
}
