package main

import (
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("missing input file")
	}

	var input_filepath = os.Args[1]
	_, err := parseInputFile(input_filepath)

	if err != nil {
		print(err.Error())
	}
}
