// .
package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("😱 : %v", err)
		}
	}()

	err := runGogeur()
	if err != nil {
		println("😱 :", err.Error())
	}
}

var errNoInputFileProvided = errors.New("No input file provided")

func runGogeur() error {
	if len(os.Args) != 2 {
		return errNoInputFileProvided
	}

	inputFilepath := os.Args[1]
	simulation, err := parseInputFile(inputFilepath)
	if err != nil {
		return err
	}
	err = verifySimulationValidity(simulation)
	if err != nil {
		return err
	}

	return err
}
