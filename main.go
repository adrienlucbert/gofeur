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

	var error = runGogeur()

	if error != nil {
		println("😱 :", error.Error())
	}
}

func runGogeur() error {
	if len(os.Args) != 2 {
		return errors.New("No input file provided")
	}

	var input_filepath = os.Args[1]
	var simulation, err = parseInputFile(input_filepath)
	if err != nil {
		return err
	}
	err = verifySimulationValidity(simulation)
	if err != nil {
		return err
	}

	return err
}
