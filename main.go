package main

import (
	"bufio"
	"fmt"
	"github.com/adrienlucbert/gofeur/pkg"
	"os"
	"time"
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
		ui.App.QueueUpdateDraw(func() {
			time.Sleep(1 * time.Second)
			if count%2 == 0 {
				ui.OutputBox.SetTitle(fmt.Sprintf("%s QUOI?", str))
			} else {
				ui.OutputBox.SetTitle(fmt.Sprintf("%s FEUR...", str))
			}
			fmt.Fprintf(ui.OutputBox, "round: %d\n", count)
		})
		count += 1

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
	for f.Scan() {
		fmt.Println(f.Text())
	}
	gofeur := pkg.Gofeur{}
	gofeur.Init()

	go ouptut(gofeur.Ui)

	gofeur.Run()
}
