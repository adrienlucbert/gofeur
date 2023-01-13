// Package pkg .
package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/adrienlucbert/gofeur/optional"
)

type (
	tokenKind       int
	parserErrorKind int
)

type tokenParser struct {
	fieldName string
	kind      tokenKind
	value     any
}

const (
	nonEmptyStringTokenKind tokenKind = iota
	unitTokenKind
	weightTokenKind
	parcelColorTokenKind
)

const (
	invalidTokenLength parserErrorKind = iota
	invalidNumberOfTokens
	invalidUnsignedInteger
	invalidCycleNumber
	invalidWeight
)

type inputError struct {
	line    uint32
	section string
	err     error
}

func (err inputError) Error() string {
	str := fmt.Sprintf("line: %d when parsing a %s", err.line, err.section)

	if err.err != nil {
		str += fmt.Sprintf(": %s", err.err.Error())
	}
	return str
}

func (kind parserErrorKind) ToString() string {
	switch kind {
	case invalidTokenLength:
		return "invalid token length"
	case invalidNumberOfTokens:
		return "invalid number of tokens"
	case invalidUnsignedInteger:
		return "invalid unsigned integer"
	case invalidWeight:
		return "invalid weight"
	case invalidCycleNumber:
		return "invalid cycle number"
	default:
		panic("Unreachable")
	}
}

type parserError interface {
	error
	Kind() parserErrorKind
}

type fieldTokenError struct {
	fieldName string
	token     string
	kind      parserErrorKind
	err       string
}

func (err fieldTokenError) Error() string {
	return fmt.Sprintf("%s for field: %s (found: '%s')", err.kind.ToString(), err.fieldName, err.token)
}

func (err fieldTokenError) Kind() parserErrorKind {
	return err.kind
}

type tokenError struct {
	kind parserErrorKind
	err  string
}

func (err tokenError) Error() string {
	return err.Kind().ToString()
}

func (err tokenError) Kind() parserErrorKind {
	return err.kind
}

type inputFileOpenError struct {
	file string
	err  error
}

func (err inputFileOpenError) Error() string {
	return fmt.Sprintf("Error while trying to open '%s': %s", err.file, err.err.Error())
}

type inputFileError struct {
	file string
	err  error
}

func (err inputFileError) Error() string {
	return fmt.Sprintf("Error in input file '%s': %s", err.file, err.err.Error())
}

// ParseInputFile parses a simulation input file. If successful, return the
// parsed Simulation.
func ParseInputFile(file string) (Simulation, error) {
	handle, err := os.Open(file)
	if err != nil {
		return Simulation{}, inputFileOpenError{err: err}
	}
	defer handle.Close()

	simul, err := parseFromReader(handle)
	if err != nil {
		err = inputFileError{file: file, err: err}
	}
	return simul, err
}

var errInvalidLine = errors.New("invalid line")

func parseFromReader(reader io.Reader) (Simulation, error) {
	scanner := bufio.NewScanner(reader)
	parser := struct {
		section                  string
		line                     uint32
		tokens                   []string
		mayBeParsedByNextSection bool
	}{section: "warehouse"}
	simul := Simulation{}
	var err parserError

	for parser.mayBeParsedByNextSection || scanner.Scan() {
		if !parser.mayBeParsedByNextSection {
			parser.line++
			line := scanner.Text()
			parser.tokens = strings.Split(line, " ")
		}
		parser.mayBeParsedByNextSection = false

		switch parser.section {
		case "warehouse":
			simul, err = parseWarehouseSection(parser.tokens)

			if err != nil {
				return Simulation{}, inputError{
					line:    parser.line,
					section: parser.section,
					err:     err,
				}
			}
			parser.section = getNextSection(parser.section).Value()

		default:
			err = parseWarehouseEntity(parser.section, parser.tokens, &simul.warehouse)
		}

		if err != nil {
			if err.Kind() == invalidNumberOfTokens {
				nextSection := getNextSection(parser.section)

				if nextSection.HasValue() {
					parser.section = nextSection.Value()
					parser.mayBeParsedByNextSection = true
				} else {
					return Simulation{}, inputError{
						line:    parser.line,
						section: parser.section,
						err:     errInvalidLine,
					}
				}
			} else {
				break
			}
		}
	}

	if err != nil {
		return Simulation{}, inputError{
			line:    parser.line,
			section: parser.section,
			err:     err,
		}
	}
	return simul, nil
}

func getNextSection(section string) optional.Optional[string] {
	switch section {
	case "warehouse":
		return optional.New("parcel")
	case "parcel":
		return optional.New("forklift")
	case "forklift":
		return optional.New("truck")
	default:
		return optional.NewEmpty[string]()
	}
}

func parseWarehouseSection(tokens []string) (Simulation, parserError) {
	simul := Simulation{}
	warehouseTokenParsers := []tokenParser{
		{
			fieldName: "width",
			kind:      unitTokenKind,
			value:     &simul.warehouse.width,
		},
		{
			fieldName: "length",
			kind:      unitTokenKind,
			value:     &simul.warehouse.length,
		},
		{
			fieldName: "cycle",
			value:     &simul.cycle,
		},
	}

	err := parseTokens(tokens, warehouseTokenParsers)
	return simul, err
}

func parseWarehouseEntity(section string, tokens []string, warehouse *warehouse) parserError {
	var err parserError

	switch section {
	case "parcel":
		var parcel parcel
		parcel, err = parseParcel(tokens)

		if err == nil {
			warehouse.parcels = append(warehouse.parcels, parcel)
		}
	case "forklift":
		var forklift forklift
		forklift, err = parseForklift(tokens)

		if err == nil {
			warehouse.forklifts = append(warehouse.forklifts, forklift)
		}
	case "truck":
		var truck truck
		truck, err = parseTruck(tokens)

		if err == nil {
			warehouse.trucks = append(warehouse.trucks, truck)
		}
	}

	return err
}

func parseParcel(tokens []string) (parcel, parserError) {
	pkg := parcel{}
	parcelTokenParsers := []tokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyStringTokenKind,
			value:     &pkg.name,
		},
		{
			fieldName: "x",
			kind:      unitTokenKind,
			value:     &pkg.X,
		},
		{
			fieldName: "y",
			kind:      unitTokenKind,
			value:     &pkg.Y,
		},
		{
			fieldName: "weight",
			kind:      parcelColorTokenKind,
			value:     &pkg.weight,
		},
	}

	err := parseTokens(tokens, parcelTokenParsers)
	return pkg, err
}

func parseForklift(tokens []string) (forklift, parserError) {
	flt := forklift{}
	forkLiftTokenParsers := []tokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyStringTokenKind,
			value:     &flt.name,
		},
		{
			fieldName: "x",
			value:     &flt.X,
		},
		{
			fieldName: "y",
			value:     &flt.Y,
		},
	}

	err := parseTokens(tokens, forkLiftTokenParsers)
	return flt, err
}

func parseTruck(tokens []string) (truck, parserError) {
	lorry := truck{}
	truckTokenParsers := []tokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyStringTokenKind,
			value:     &lorry.name,
		},
		{
			fieldName: "x",
			kind:      unitTokenKind,
			value:     &lorry.X,
		},
		{
			fieldName: "y",
			kind:      unitTokenKind,
			value:     &lorry.Y,
		},
		{
			fieldName: "maximum_weight",
			kind:      weightTokenKind,
			value:     &lorry.maxWeight,
		},
		{
			fieldName: "available",
			value:     &lorry.available,
		},
	}

	err := parseTokens(tokens, truckTokenParsers)
	return lorry, err
}

func parseTokens(tokens []string, tokenParsers []tokenParser) parserError {
	if len(tokens) != len(tokenParsers) {
		return fieldTokenError{kind: invalidNumberOfTokens}
	}

	for i, tokenParser := range tokenParsers {
		token := tokens[i]
		maybeTokenErr := parseToken(token, tokenParser.kind, tokenParser.value)

		tokenErr, ok := maybeTokenErr.(tokenError)

		if ok {
			return fieldTokenError{
				kind:      tokenErr.kind,
				err:       tokenErr.Error(),
				fieldName: tokenParser.fieldName,
				token:     token,
			}
		}
	}
	return nil
}

func parseToken(token string, kind tokenKind, value any) error {
	var err error

	switch ptr := value.(type) {
	case *string:
		err = parseStringToken(token, kind, ptr)
	case *uint32:
		err = parseUint32Token(token, kind, ptr)
	case *simulationCycle:
		err = parseSimulationCycleToken(token, kind, ptr)
	case *gridUnit:
		err = parseUnitToken(token, kind, ptr)
	case *weight:
		err = parseWeightToken(token, kind, ptr)
	default:
		panic("Unreachable: Unexpected pointer type")
	}

	return err
}

func parseStringToken(token string, kind tokenKind, ptr *string) parserError {
	switch kind {
	case nonEmptyStringTokenKind:
		if len(token) > 0 {
			*ptr = token
		} else {
			return tokenError{kind: invalidTokenLength, err: "is an empty string"}
		}
	default:
		panic("Unreachable")
	}
	return nil
}

func parseUint32Token(token string, _ tokenKind, ptr *uint32) parserError {
	value, err := parseUint32Field(token)

	if err == nil {
		*ptr = value
		return nil
	}
	return tokenError{kind: invalidUnsignedInteger, err: err.Error()}
}

func parseSimulationCycleToken(token string, _ tokenKind, ptr *simulationCycle) parserError {
	value, err := parseUint32Field(token)

	if err == nil {
		if value < 10 || value > 100000 {
			return tokenError{kind: invalidCycleNumber, err: "should be between 10 and 100_000"}
		}
		*ptr = simulationCycle(value)
		return nil
	}
	return tokenError{kind: invalidUnsignedInteger, err: err.Error()}
}

func parseUnitToken(token string, kind tokenKind, ptr *gridUnit) parserError {
	var value uint32
	err := parseUint32Token(token, kind, &value)

	if err == nil {
		*ptr = gridUnit(value)
	}
	return err
}

func parseWeightToken(token string, kind tokenKind, ptr *weight) parserError {
	switch kind {
	case weightTokenKind:
		value, err := parseUint32Field(token)

		if err == nil {
			*ptr = weight(value)
		} else {
			return tokenError{kind: invalidWeight, err: err.Error()}
		}
	case parcelColorTokenKind:
		parcelWeight, err := parseWeight(token)

		if err == nil {
			*ptr = parcelWeight
		} else {
			return tokenError{kind: invalidWeight, err: err.Error()}
		}
	default:
		panic("Unreachable: Unexpected TokenParserKind")
	}

	return nil
}

func parseUint32Field(token string) (uint32, error) {
	value, err := strconv.ParseUint(token, 10, 32)
	return uint32(value), err
}

var errInvalidColor = errors.New("invalid color")

func parseWeight(maybeColor string) (weight, error) {
	colors := []string{"yellow", "green", "blue"}
	colorsWeight := []weight{yellow, green, blue}

	for i, color := range colors {
		if color == maybeColor || strings.ToUpper(color) == maybeColor {
			return colorsWeight[i], nil
		}
	}

	return yellow, errInvalidColor
}
