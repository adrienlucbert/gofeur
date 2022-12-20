package main

import (
	"bufio"
	"fmt"
	"github.com/adrienlucbert/gofeur/pkg"
	"os"
)

func getFileContent(file string) (*os.File, *bufio.Scanner) {
	fd, err := os.Open(file)
	if err != nil {
		println(err)
		panic("cannot open file")
	}
	fmt.Println(file)
	fileScanner := bufio.NewScanner(fd)
	fileScanner.Split(bufio.ScanLines)
	return fd, fileScanner
}

func main() {
	fmt.Println(os.Args)

	if len(os.Args) != 2 {
		panic("missing input file")
	}
	file := os.Args[1]
	fd, f := getFileContent(file)
	defer fd.Close()
	for f.Scan() {
		fmt.Println(f.Text())
	}
	ui := ui.Start()
	if err := ui.App.SetRoot(ui.FlexLayout, true).SetFocus(ui.FlexLayout).Run(); err != nil {
		panic(err)
	}
}
