// main function
package main

import (
	"bufio"
	"os"

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
	if len(os.Args) != 2 {
		panic("missing input file")
	}
	file := os.Args[1]
	fd, f := getFileContent(file)
	defer fd.Close()

	gofeur := pkg.ParseFile(f)

	gofeur.Init()
	go pkg.Logic(gofeur.Ui)
	gofeur.Run()
}
