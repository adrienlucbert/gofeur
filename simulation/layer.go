package simulation

import (
	"time"
)

// Layer is the application layer responsible for managing the game logic
type Layer struct {
	Simulation *Simulation
}

// Attach initializes the LogicLayer
func (layer *Layer) Attach() {
	layer.Simulation.start()
}

// Update runs the game logic
func (layer *Layer) Update(elapsedTime time.Duration) {
	layer.Simulation.simulateRound()
}

// Detach handles the game end
func (layer *Layer) Detach() {
	layer.Simulation.terminate()
	// TODO: print end-of-game data
}
