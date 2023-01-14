package pkg

import (
	"fmt"
	"math"
)

// Vector struct stores 2 ints, defining a position or a movement
type Vector struct {
	X int
	Y int
}

func (v Vector) String() string {
	return fmt.Sprintf("(%d;%d)", v.X, v.Y)
}

// Add returns a vector resulting from the addition of 2 vectors
func (v Vector) Add(rhs Vector) Vector {
	v.X += rhs.X
	v.Y += rhs.Y
	return v
}

// SquaredDistance calculates the squared distance between 2 vectors
func (v Vector) SquaredDistance(rhs Vector) float32 {
	return float32(math.Pow(float64(v.X-rhs.X), 2) + math.Pow(float64(v.Y-rhs.Y), 2))
}
