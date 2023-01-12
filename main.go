// .
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/adrienlucbert/gofeur/pkg"
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

type parserError struct {
	parser error
}

func (err parserError) Error() string {
	return fmt.Sprintf("Parser: %s", err.parser.Error())
}

type simulationValidityError struct {
	validation error
}

func (err simulationValidityError) Error() string {
	return fmt.Sprintf("Simulation validity: %s", err.validation.Error())
}

func runGogeur() error {
	if len(os.Args) != 2 {
		return errNoInputFileProvided
	}

	inputFilepath := os.Args[1]
	simulation, err := pkg.ParseInputFile(inputFilepath)
	if err != nil {
		return parserError{parser: err}
	}
	err = pkg.VerifySimulationValidity(simulation)
	if err != nil {
		return simulationValidityError{validation: err}
	}

	return nil
}
