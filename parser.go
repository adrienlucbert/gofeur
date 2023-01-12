package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adrienlucbert/gofeur/optional"
)

type TokenKind int
type ParserErrorKind int

type TokenParser struct {
	fieldName string
	kind      TokenKind
	value     any
}

const (
	nonEmptyString TokenKind = iota
	unit
	weight
	parcelColor
)

const (
	invalidTokenLength ParserErrorKind = iota
	invalidNumberOfTokens
	invalidUnsignedInteger
	invalidWeight
)

type InputFileError struct {
	filename string
	line     uint32
	section  string
	err      error
}

func (err InputFileError) Error() string {
	var str = fmt.Sprintf("Error in input file '%s' (line: %d) when parsing a %s", err.filename, err.line, err.section)

	if err.err != nil {
		str = str + fmt.Sprintf(": %s", err.err.Error())
	}
	return str
}

func (kind ParserErrorKind) ToString() string {
	switch kind {
	case invalidTokenLength:
		return "invalid token length"
	case invalidNumberOfTokens:
		return "invalid number of tokens"
	case invalidUnsignedInteger:
		return "invalid unsigned integer"
	case invalidWeight:
		return "invalid weight"
	default:
		panic("Unreachable")
	}
}

type ParserError interface {
	error
	Kind() ParserErrorKind
}

type FieldTokenError struct {
	fieldName string
	token     string
	kind      ParserErrorKind
	err       string
}

func (err FieldTokenError) Error() string {
	return fmt.Sprintf("%s for field: %s (found: '%s')", err.kind.ToString(), err.fieldName, err.token)
}

func (err FieldTokenError) Kind() ParserErrorKind {
	return err.kind
}

type TokenError struct {
	kind ParserErrorKind
	err  string
}

func (err TokenError) Error() string {
	return fmt.Sprintf("%s", err.Kind().ToString())
}

func (err TokenError) Kind() ParserErrorKind {
	return err.kind
}

type parseParcelError struct {
	kind  ParserErrorKind
	token string
}

func parseInputFile(file string) (Simulation, error) {
	fd, f := getFileContent(file)
	defer fd.Close()

	var parser = struct {
		section                       string
		line                          uint32
		tokens                        []string
		may_be_parsed_by_next_section bool
	}{section: "warehouse"}
	var simulation = Simulation{}
	var err ParserError

	for parser.may_be_parsed_by_next_section || f.Scan() {
		if err == nil && !parser.may_be_parsed_by_next_section {
			parser.line += 1
			var line = f.Text()
			parser.tokens = strings.Split(line, " ")
		}
		parser.may_be_parsed_by_next_section = false

		switch parser.section {
		case "warehouse":
			simulation, err = parseWarehouseSection(parser.tokens)

			if err != nil {
				return simulation, InputFileError{
					filename: file,
					line:     parser.line,
					section:  parser.section,
					err:      err,
				}

			}
			parser.section = getNextSection(parser.section).Value()
		case "parcel":
			var parcel = Parcel{}
			parcel, err = parseParcel(parser.tokens)

			if err == nil {
				simulation.warehouse.parcels = append(simulation.warehouse.parcels, parcel)
			}
		case "forklift":
			var forklift = Forklift{}
			forklift, err = parseForklift(parser.tokens)

			if err == nil {
				simulation.warehouse.forklifts = append(simulation.warehouse.forklifts, forklift)
			}
		case "truck":
			var truck = Truck{}
			truck, err = parseTruck(parser.tokens)

			if err == nil {
				simulation.warehouse.trucks = append(simulation.warehouse.trucks, truck)
			}
		}

		if err != nil {
			if err.Kind() == invalidNumberOfTokens {
				var next_section = getNextSection(parser.section)

				if next_section.HasValue() {
					parser.section = next_section.Value()
					parser.may_be_parsed_by_next_section = true
				} else {
					return Simulation{}, InputFileError{
						filename: file,
						line:     parser.line,
						section:  parser.section,
						err:      errors.New("invalid line"),
					}
				}
			} else {
				break
			}
		}
	}

	if err != nil {
		return simulation, InputFileError{
			filename: file,
			line:     parser.line,
			section:  parser.section,
			err:      err,
		}
	}
	return simulation, err
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

func parseWarehouseSection(tokens []string) (Simulation, ParserError) {
	var simulation = Simulation{}
	var warehouseTokenParsers = []TokenParser{
		{
			fieldName: "width",
			kind:      unit,
			value:     &simulation.warehouse.width,
		},
		{
			fieldName: "length",
			kind:      unit,
			value:     &simulation.warehouse.length,
		},
		{
			fieldName: "length",
			value:     &simulation.cycle,
		},
	}

	var err = parseTokens(tokens, warehouseTokenParsers)
	return simulation, err
}

func parseParcel(tokens []string) (Parcel, ParserError) {
	var parcel = Parcel{}
	var parcelTokenParsers = []TokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyString,
			value:     &parcel.name,
		},
		{
			fieldName: "x",
			kind:      unit,
			value:     &parcel.X,
		},
		{
			fieldName: "y",
			kind:      unit,
			value:     &parcel.Y,
		},
		{
			fieldName: "weight",
			kind:      parcelColor,
			value:     &parcel.weight,
		},
	}

	var err = parseTokens(tokens, parcelTokenParsers)
	return parcel, err
}

func parseForklift(tokens []string) (Forklift, ParserError) {
	var forklift = Forklift{}
	var forkLiftTokenParsers = []TokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyString,
			value:     &forklift.name,
		},
		{
			fieldName: "x",
			value:     &forklift.X,
		},
		{
			fieldName: "y",
			value:     &forklift.Y,
		},
	}

	var err = parseTokens(tokens, forkLiftTokenParsers)
	return forklift, err
}

func parseTruck(tokens []string) (Truck, ParserError) {
	var truck = Truck{}
	var truckTokenParsers = []TokenParser{
		{
			fieldName: "name",
			kind:      nonEmptyString,
			value:     &truck.name,
		},
		{
			fieldName: "x",
			kind:      unit,
			value:     &truck.X,
		},
		{
			fieldName: "y",
			kind:      unit,
			value:     &truck.Y,
		},
		{
			fieldName: "maximum_weight",
			kind:      weight,
			value:     &truck.max_weight,
		},
		{
			fieldName: "available",
			value:     &truck.available,
		},
	}

	var err = parseTokens(tokens, truckTokenParsers)
	return truck, err
}

func parseTokens(tokens []string, tokenParsers []TokenParser) ParserError {
	if len(tokens) != len(tokenParsers) {
		return FieldTokenError{kind: invalidNumberOfTokens}
	}

	for i, tokenParser := range tokenParsers {
		var token = tokens[i]
		var maybe_token_err = parseToken(token, tokenParser.kind, tokenParser.value)

		token_err, ok := maybe_token_err.(TokenError)

		if ok {
			return FieldTokenError{
				kind:      token_err.kind,
				err:       token_err.Error(),
				fieldName: tokenParser.fieldName,
				token:     token,
			}
		}
	}
	return nil
}

func parseToken(token string, kind TokenKind, value any) error {
	var err error

	switch ptr := value.(type) {
	case *string:
		err = parseStringToken(token, kind, ptr)
	case *uint32:
		err = parseUint32Token(token, kind, ptr)
	case *Unit:
		err = parseUnitToken(token, kind, ptr)
	case *Weight:
		err = parseWeightToken(token, kind, ptr)
	default:
		panic("Unreachable: Unexpected pointer type")
	}

	return err
}

func parseStringToken(token string, kind TokenKind, ptr *string) ParserError {
	switch kind {
	case nonEmptyString:
		if len(token) > 0 {
			*ptr = token
		} else {
			return TokenError{kind: invalidTokenLength, err: "is an empty string"}
		}
	default:
		panic("Unreachable")
	}
	return nil
}

func parseUint32Token(token string, kind TokenKind, ptr *uint32) ParserError {
	var value, err = parseUint32Field(token)

	if err == nil {
		*ptr = value
		return nil
	}
	return TokenError{kind: invalidUnsignedInteger, err: err.Error()}
}

func parseUnitToken(token string, kind TokenKind, ptr *Unit) ParserError {
	var value uint32
	var err = parseUint32Token(token, kind, &value)

	if err == nil {
		*ptr = (Unit)(value)
	}
	return err
}

func parseWeightToken(token string, kind TokenKind, ptr *Weight) ParserError {
	switch kind {
	case weight:
		var value, err = parseUint32Field(token)

		if err == nil {
			*ptr = (Weight)(value)
		} else {
			return TokenError{kind: invalidWeight, err: err.Error()}
		}
	case parcelColor:
		var parcelWeight, err = parseWeight(token)

		if err == nil {
			*ptr = parcelWeight
		} else {
			return TokenError{kind: invalidWeight, err: err.Error()}
		}
	default:
		panic("Unreachable: Unexpected TokenParserKind")
	}

	return nil
}

func parseUint32Field(token string) (uint32, error) {
	var value, err = strconv.ParseUint(token, 10, 32)
	return (uint32)(value), err
}

func parseWeight(maybe_color string) (Weight, error) {
	var colors = []string{"yellow", "green", "blue"}
	var colors_weight = []Weight{yellow, green, blue}

	for i, color := range colors {
		if color == maybe_color || strings.ToUpper(color) == maybe_color {
			return colors_weight[i], nil
		}
	}

	return yellow, errors.New("invalid color")
}
