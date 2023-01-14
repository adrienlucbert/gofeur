// Package pkg provides common utils
package pkg

import "time"

// Layer represents an application layer that is updated on every update of the
// application main loop.
type Layer interface {
	// Attach is called once when initializing the layer
	Attach()
	// Update is called after every game update
	Update(elapsedTime time.Duration)
	// Detach is called once when the game is over
	Detach()
}
