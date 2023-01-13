package parsing

import (
	"bufio"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	STARTUP = iota
	PARCEL
	FORKLIFT
	TRUCK
)

func colorToUint8(color string) uint8 {
	colorLower := strings.ToLower(color)

	if colorLower == "green" {
		return uint8(Green)
	} else if colorLower == "yellow" {
		return uint8(Yellow)
	}
	return uint8(Blue)
}

func factory[T any](elem T, args ...string) {
	va := reflect.ValueOf(elem)

	for i, arg := range args {
		field := va.Elem().Field(i)
		fmt.Printf("\ttype %s - %d\n", field.Type().Name(), i)
		fieldType := field.Type().Name()
		switch fieldType {
		case "uint":
			num, err := strconv.ParseUint(arg, 10, 64)
			if err == nil {
				field.SetUint(num)
			} else {
				panic("Error: factory, cannot SetUint")
			}
		case "Color":
			field.SetUint(uint64(colorToUint8(arg)))
		default:
			field.SetString(arg)
		}
	}
}

func parseLine(index int, line string, gofeur *Gofeur) {
	splitted := strings.Fields(line)

	switch index {
	case STARTUP:
		if (Startup{}) != gofeur.ST {
			fmt.Println(`Startup values has already been initialized, 
            the line will be ignored`)
			return
		}
		factory(&gofeur.ST, splitted...)
	case PARCEL:
		newParcel := Parcel{}
		factory(&newParcel, splitted...)
		gofeur.SB.Packs = append(gofeur.SB.Packs, newParcel)
	case FORKLIFT:
		newForklift := Forklift{}
		factory(&newForklift, splitted...)
		gofeur.SB.Forklifts = append(gofeur.SB.Forklifts, newForklift)
	case TRUCK:
		newTruck := Truck{}
		factory(&newTruck, splitted...)
		gofeur.SB.Trucks = append(gofeur.SB.Trucks, newTruck)
	}
}

func ParseFile(file *bufio.Scanner) *Gofeur {
	gofeur := Gofeur{}
	patterns := []string{
		`^\d+\s+\d+\s+\d+$`,
		`^[A-Za-zÀ-ÿ0-9\_]+\s+\d+\s+\d+\s+(?i)green|blue|yellow$`,
		`^[A-Za-zÀ-ÿ0-9\_]+\s+\d+\s+\d+$`,
		`^[A-Za-zÀ-ÿ0-9\_]+\s+\d+\s+\d+\s+\d+\s+\d+$`,
	}
	isValid := false

	for file.Scan() {
		line := file.Text()
		for i, pattern := range patterns {
			matched, _ := regexp.MatchString(pattern, line)
			if matched {
				fmt.Printf("matched ->[%s]<- with -> %s\n", line, pattern)
				parseLine(i, line, &gofeur)
				isValid = matched
				break
			}
		}
		if !isValid {
			fmt.Printf("Error this line is not valid: %s\n", line)
			return nil
		}
		isValid = false
	}
	return &gofeur
}
