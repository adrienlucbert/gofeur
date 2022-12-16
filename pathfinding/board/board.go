// Package board contains types and utils related to a pathfinding board
package board

// Tile holds information useful to
type Tile struct {
	Blocked bool
}

func (t Tile) String() string {
	if t.Blocked {
		return "#"
	}
	return "·"
}

// Board a 2D array of tiles by the pathfinding module
type Board [][]Tile

// Width returns the board's width
func (b *Board) Width() int {
	if b.Height() == 0 {
		return 0
	}
	return len((*b)[0])
}

// Height returns the board's height
func (b *Board) Height() int {
	return len(*b)
}

func (b *Board) String() string {
	s := ""
	for _, row := range *b {
		for _, col := range row {
			s += col.String() + " "
		}
		s += "\n"
	}
	return s
}

// Returns the tile at given coordinates
func (b *Board) At(x int, y int) *Tile {
	return &(*b)[y][x]
}

// Contains returns whether the given position is within the board's bounds
func (b *Board) Contains(x int, y int) bool {
	switch {
	case x < 0 || x > b.Width():
		return false
	case y < 0 || y > b.Height():
		return false
	default:
		return true
	}
}

// New initializes a board of the given size
func New(width uint, height uint) Board {
	board := make(Board, height)
	for i := range board {
		board[i] = make([]Tile, width)
	}
	return board
}
