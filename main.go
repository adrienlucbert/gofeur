// main function
package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/adrienlucbert/gofeur/pkg"
)

func getFileContent(file string) (*os.File, *bufio.Scanner) {
	fd, err := os.Open(file)
	if err != nil {
		println(err)
		panic("cannot open file")
	}
	fileScanner := bufio.NewScanner(fd)
	fileScanner.Split(bufio.ScanLines)
	return fd, fileScanner
}

func main() {
	filename := flag.String("filename", "", "Map file path")
	displayUI := flag.Bool("ui", false, "Display UI")
	flag.Parse()
	if *filename == "" {
		panic("missing input file")
	}
	fd, f := getFileContent(*filename)
	defer fd.Close()

	gofeur := pkg.ParseFile(f)
	simulation := pkg.NewSimulation(gofeur)

	layers := []pkg.Layer{
		&pkg.LogicLayer{Simulation: &simulation},
	}
	if *displayUI {
		layers = append(layers, &pkg.UILayer{Gofeur: gofeur, Simulation: &simulation})
	}
	for _, layer := range layers {
		layer.Attach()
	}
	lastUpdateTime := time.Now()
	for simulation.Status == pkg.Running {
		updateTime := time.Now()
		elapsedTime := updateTime.Sub(lastUpdateTime)
		lastUpdateTime = updateTime
		for _, layer := range layers {
			layer.Update(elapsedTime)
		}
	}
	for _, layer := range layers {
		layer.Detach()
	}
}
