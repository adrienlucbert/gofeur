package pathfinding

import (
	"testing"

	"github.com/adrienlucbert/gofeur/pathfinding/board"
	"github.com/stretchr/testify/assert"
)

func TestPathfindingFindsShortestPath1(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	path, err := Resolve(&b, Vector{1, 0}, Vector{4, 3})
	assert.Nil(t, err)
	shortestPath := []Vector{{1, 1}, {1, 2}, {2, 2}, {3, 2}, {3, 3}, {4, 3}}
	assert.Equal(t, path, shortestPath)
}

func TestPathfindingFindsShortestPath2(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	b.At(1, 2).Blocked = true
	path, err := Resolve(&b, Vector{1, 0}, Vector{4, 3})
	assert.Nil(t, err)
	shortestPath := []Vector{{2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}}
	assert.Equal(t, path, shortestPath)
}

func TestPathfindingNoPath(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 0).Blocked = true
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	b.At(1, 2).Blocked = true
	path, err := Resolve(&b, Vector{1, 0}, Vector{4, 3})
	assert.Equal(t, err, ErrPathNotFound)
	assert.Empty(t, path)
}
