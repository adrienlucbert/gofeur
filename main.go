package main

import (
	"bufio"
	"fmt"
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
	fmt.Println(file)
	fileScanner := bufio.NewScanner(fd)
	fileScanner.Split(bufio.ScanLines)
	return fd, fileScanner
}

func ouptut(ui *pkg.UI) {
	str := "Go "
	count := 0
	for {
		time.Sleep(4 * time.Second)
		if count%2 == 0 {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s QUOI?", str))
		} else {
			ui.OutputBox.SetTitle(fmt.Sprintf("%s FEUR...", str))
		}
		fmt.Fprintf(ui.OutputBox, "round: %d\n", count)
		count++
		ui.App.Draw()
	}
}

func main() {
	fmt.Println(os.Args)

	if len(os.Args) != 2 {
		panic("missing input file")
	}
	file := os.Args[1]
	fd, f := getFileContent(file)
	defer fd.Close()

    gofeur := pkg.ParseFile(f)

	gofeur.Init()
	go ouptut(gofeur.Ui)
	gofeur.Run()
}
