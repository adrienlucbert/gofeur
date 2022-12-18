// Package pathfinding provides pathfinding utils
package pathfinding

import (
	"errors"
	"fmt"
	"math"

	"github.com/adrienlucbert/gofeur/optional"
	"github.com/adrienlucbert/gofeur/pathfinding/board"
)

// Vector struct stores 2 ints, defining a position or a movement
type Vector struct {
	X int
	Y int
}

func (v Vector) String() string {
	return fmt.Sprintf("(%d;%d)", v.X, v.Y)
}

func (v Vector) add(rhs Vector) Vector {
	v.X += rhs.X
	v.Y += rhs.Y
	return v
}

func (v Vector) squaredDistance(rhs Vector) float32 {
	return float32(math.Pow(float64(v.X-rhs.X), 2) + math.Pow(float64(v.Y-rhs.Y), 2))
}

type node struct {
	parent   *node
	position Vector
	g        float32
	h        float32
	f        float32
}

func findBestNodeInList(list []node) (int, node) {
	bestNode := optional.NewEmpty[node]()
	bestIndex := optional.NewEmpty[int]()
	for index, node := range list {
		if !bestNode.HasValue() || node.f < bestNode.Value().f {
			bestIndex.Set(index)
			bestNode.Set(list[index])
		}
	}
	return bestIndex.Value(), bestNode.Value()
}

func findNodeInList(n node, list []node) optional.Optional[node] {
	for _, it := range list {
		if it.position == n.position {
			return optional.New(it)
		}
	}
	return optional.NewEmpty[node]()
}

func reconstructPath(current *node) []Vector {
	path := []Vector{}
	for current != nil && current.parent != nil {
		// PERF: reverse list once at the end instead of prepending each node
		path = append([]Vector{current.position}, path...)
		current = current.parent
	}
	return path
}

func distanceBasedHeuristic(current Vector, end Vector) float32 {
	return current.squaredDistance(end)
}

// ErrPathNotFound is returned when a path couldn't be found
var ErrPathNotFound = errors.New("Couldn't find a path")

func astar(b *board.Board, start Vector, end Vector, heuristic func(Vector) float32) ([]Vector, error) {
	openQueue := []node{}
	closedQueue := []node{}
	openQueue = append(openQueue, node{position: start})
	for len(openQueue) > 0 {
		bestNodeIndex, bestNode := findBestNodeInList(openQueue)

		openQueue = append(openQueue[:bestNodeIndex], openQueue[bestNodeIndex+1:]...)
		closedQueue = append(closedQueue, bestNode)

		if bestNode.position == end {
			// Reached the end
			return reconstructPath(&bestNode), nil
		}

		children := []node{}
		for _, direction := range []Vector{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			child := node{
				parent:   &bestNode,
				position: bestNode.position.add(direction),
			}
			if !b.Contains(child.position.X, child.position.Y) {
				continue
			}
			if b.At(child.position.X, child.position.Y).Blocked {
				continue
			}
			children = append(children, child)
		}

		for _, child := range children {
			if findNodeInList(child, closedQueue).HasValue() {
				continue
			}
			child.g = bestNode.g + 1
			child.h = heuristic(child.position)
			child.f = child.g + child.h
			openNode := findNodeInList(child, openQueue)
			if openNode.HasValue() && child.g > openNode.Value().g {
				continue
			}
			openQueue = append(openQueue, child)
		}
	}
	return []Vector{}, ErrPathNotFound
}

// ResolveH returns, if possible, a series of moves that form a path between
// start and end, using a provided heuristic
func ResolveH(maze *board.Board, start Vector, end Vector, heuristic func(Vector) float32) ([]Vector, error) {
	return astar(maze, start, end, heuristic)
}

// Resolve returns, if possible, a series of moves that form a path between
// start and end, using a distance-based heuristic
func Resolve(maze *board.Board, start Vector, end Vector) ([]Vector, error) {
	heuristic := func(n Vector) float32 {
		return distanceBasedHeuristic(n, end)
	}
	return astar(maze, start, end, heuristic)
}
