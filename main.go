package main

import (
	"fmt"
	"github.com/adrienlucbert/gofeur/pathfinding/board"
)

func main() {
	b := board.New(5, 4)
	b.At(2, 1).Blocked = true
	b.At(1, 3).Blocked = true
	fmt.Printf("%d x %d\n", b.Width(), b.Height())
	fmt.Println(b.String())
}
