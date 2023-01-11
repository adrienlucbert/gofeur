// main function
package main

import (
	"bufio"
	"flag"
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
	fileScanner := bufio.NewScanner(fd)
	fileScanner.Split(bufio.ScanLines)
	return fd, fileScanner
}

func main() {
	filename := flag.String("filename", "", "Map file path")
	flag.Parse()
	if *filename == "" {
		panic("missing input file")
	}
	fd, f := getFileContent(*filename)
	defer fd.Close()

	gofeur := pkg.ParseFile(f)

	layers := []pkg.Layer{}
	for _, layer := range layers {
		layer.Attach()
	}
	lastUpdateTime := time.Now()
	for gofeur.Status == pkg.Running {
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
