package pathfinding

import (
	"testing"

	"github.com/adrienlucbert/gofeur/pathfinding/board"
	"github.com/adrienlucbert/gofeur/pkg"
	"github.com/stretchr/testify/assert"
)

func TestPathfindingFindsShortestPath1(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	path, err := Resolve(&b, pkg.Vector{X: 1, Y: 0}, pkg.Vector{X: 4, Y: 3})
	assert.Nil(t, err)
	shortestPath := []pkg.Vector{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 3}}
	assert.Equal(t, path, shortestPath)
}

func TestPathfindingFindsShortestPath2(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	b.At(1, 2).Blocked = true
	path, err := Resolve(&b, pkg.Vector{X: 1, Y: 0}, pkg.Vector{X: 4, Y: 3})
	assert.Nil(t, err)
	shortestPath := []pkg.Vector{{X: 2, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 3}}
	assert.Equal(t, path, shortestPath)
}

func TestPathfindingNoPath(t *testing.T) {
	b := board.New(5, 4)
	b.At(2, 0).Blocked = true
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	b.At(1, 2).Blocked = true
	path, err := Resolve(&b, pkg.Vector{X: 1, Y: 0}, pkg.Vector{X: 4, Y: 3})
	assert.Equal(t, err, ErrPathNotFound)
	assert.Empty(t, path)
}
