// Package pathfinding provides pathfinding utils
package pathfinding

import (
	"errors"

	"github.com/adrienlucbert/gofeur/optional"
	"github.com/adrienlucbert/gofeur/pathfinding/board"
	"github.com/adrienlucbert/gofeur/pkg"
)

// ResolveH returns, if possible, a series of moves that form a path between
// start and end, using a provided heuristic
func ResolveH(maze *board.Board, start pkg.Vector, end pkg.Vector, heuristic func(pkg.Vector) float32) ([]pkg.Vector, error) {
	return astar(maze, start, end, heuristic)
}

func distanceBasedHeuristic(current pkg.Vector, end pkg.Vector) float32 {
	return current.SquaredDistance(end)
}

// Resolve returns, if possible, a series of moves that form a path between
// start and end, using a distance-based heuristic
func Resolve(maze *board.Board, start pkg.Vector, end pkg.Vector) ([]pkg.Vector, error) {
	heuristic := func(n pkg.Vector) float32 {
		return distanceBasedHeuristic(n, end)
	}
	return astar(maze, start, end, heuristic)
}

type node struct {
	parent   *node
	position pkg.Vector
	g        float32
	h        float32
	f        float32
}

// ErrPathNotFound is returned when a path couldn't be found
var ErrPathNotFound = errors.New("Couldn't find a path")

func astar(b *board.Board, start pkg.Vector, end pkg.Vector, heuristic func(pkg.Vector) float32) ([]pkg.Vector, error) {
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
		for _, direction := range []pkg.Vector{{X: 0, Y: 1}, {X: 0, Y: -1}, {X: 1, Y: 0}, {X: -1, Y: 0}} {
			child := node{
				parent:   &bestNode,
				position: bestNode.position.Add(direction),
			}
			if isPositionAvailable(b, child.position) {
				children = append(children, child)
			}
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
	return []pkg.Vector{}, ErrPathNotFound
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

func isPositionAvailable(maze *board.Board, position pkg.Vector) bool {
	if position.X < 0 || position.Y < 0 || !maze.IsInBounds(uint(position.X), uint(position.Y)) {
		return false
	}
	if maze.At(uint(position.X), uint(position.Y)).Blocked {
		return false
	}
	return true
}

func reconstructPath(current *node) []pkg.Vector {
	path := []pkg.Vector{}
	for current != nil && current.parent != nil {
		// PERF: reverse list once at the end instead of prepending each node
		path = append([]pkg.Vector{current.position}, path...)
		current = current.parent
	}
	return path
}

func findNodeInList(n node, list []node) optional.Optional[node] {
	for _, it := range list {
		if it.position == n.position {
			return optional.New(it)
		}
	}
	return optional.NewEmpty[node]()
}
