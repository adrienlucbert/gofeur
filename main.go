// main function
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/adrienlucbert/gofeur/config"
	"github.com/adrienlucbert/gofeur/logger"
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
	"github.com/adrienlucbert/gofeur/simulation"
	"github.com/adrienlucbert/gofeur/ui"
)

type gofeurError struct {
	err string
}

func (err gofeurError) Error() string {
	return fmt.Sprintf("😱 : %s", err.err)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			var gofeurError = gofeurError{err: fmt.Sprintf("%v", err)}

			println(gofeurError.Error())
		}
	}()

	filename := flag.String("filename", "", "Map file path")
	displayUI := flag.Bool("ui", false, "Display UI")
	logLevel := flag.String("log-level", "Debug", "Log level (Debug, Info, Warn, Error, None)")
	flag.Parse()

	config.Set("displayUI", displayUI)
	logger.SetLogLevel(*logLevel)

	gofeur, err := parsing.ParseInputFile(*filename)
	if err != nil {
		println(gofeurError{err: err.Error()}.Error())
		return
	}
	err = parsing.VerifySimulationValidity(gofeur)
	if err != nil {
		println(gofeurError{err: err.Error()}.Error())
		return
	}

	sim := simulation.New(&gofeur)

	layers := []pkg.Layer{
		&simulation.Layer{Simulation: &sim},
	}
	if *displayUI {
		layers = append(layers, &ui.Layer{Gofeur: &gofeur, Simulation: &sim})
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
