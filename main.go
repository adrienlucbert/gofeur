// main function
package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/adrienlucbert/gofeur/config"
	"github.com/adrienlucbert/gofeur/logger"
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
	"github.com/adrienlucbert/gofeur/simulation"
	"github.com/adrienlucbert/gofeur/ui"
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
	logLevel := flag.String("log-level", "Debug", "Log level (Debug, Info, Warn, Error, None)")
	flag.Parse()

	config.Set("displayUI", displayUI)
	logger.SetLogLevel(*logLevel)

	if *filename == "" {
		panic("missing input file")
	}
	fd, f := getFileContent(*filename)
	defer fd.Close()

	gofeur := parsing.ParseFile(f)
	sim := simulation.New(gofeur)

	layers := []pkg.Layer{
		&simulation.Layer{Simulation: &sim},
	}
	if *displayUI {
		layers = append(layers, &ui.Layer{Gofeur: gofeur, Simulation: &sim})
	}
	for _, layer := range layers {
		layer.Attach()
	}
	lastUpdateTime := time.Now()
	for sim.IsRunning() {
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
