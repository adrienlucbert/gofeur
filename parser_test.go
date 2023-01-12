package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseReader(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Simulation
		hasError       bool
	}

	var testCases = []testCase{
		{
			input: []string{
				"5 4 1000",
				"colis_a_livrer 2 1 green",
				"paquet 2 2 BLUE",
				"deadpool 0 3 yellow",
				"colère_DU_dragon 4 1 green",
				"transpalette_1 0 5",
				"camion_b 3 4 4000 5",
				"camion_a 2 2 4007 4",
			},
			expectedOutput: Simulation{
				cycle: 1000,
				warehouse: Warehouse{
					width:  5,
					length: 4,
					parcels: []Parcel{
						{
							name:       "colis_a_livrer",
							Coordinate: Coordinate{X: 2, Y: 1},
							weight:     green,
						},
						{
							name:       "paquet",
							Coordinate: Coordinate{X: 2, Y: 2},
							weight:     blue,
						},
						{
							name:       "deadpool",
							Coordinate: Coordinate{X: 0, Y: 3},
							weight:     yellow,
						},
						{
							name: "colère_DU_dragon",
							Coordinate: Coordinate{
								X: 4,
								Y: 1,
							},
							weight: green,
						},
					},
					forklifts: []Forklift{
						{
							name:       "transpalette_1",
							Coordinate: Coordinate{X: 0, Y: 5},
						},
					},
					trucks: []Truck{
						{
							name:       "camion_b",
							Coordinate: Coordinate{X: 3, Y: 4},
							max_weight: 4000,
							available:  5,
						},
						{
							name:       "camion_a",
							Coordinate: Coordinate{X: 2, Y: 2},
							max_weight: 4007,
							available:  4,
						},
					},
				},
			},
		},
		{
			input: []string{
				"10 50 243",
				"forklift 1 10",
				"truck 0 5 10000 60",
			},
			expectedOutput: Simulation{
				cycle: 243,
				warehouse: Warehouse{
					width: 10, length: 50,
					forklifts: []Forklift{
						{
							name:       "forklift",
							Coordinate: Coordinate{X: 1, Y: 10},
						},
					},
					trucks: []Truck{
						{
							name:       "truck",
							Coordinate: Coordinate{X: 0, Y: 5},
							max_weight: 10000,
							available:  60,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		var input = strings.Join(testCase.input, "\n")
		var reader = strings.NewReader(input)
		var simulation, err = parseFromReader(reader)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, simulation, testCase.expectedOutput)
		}
	}
}

func TestParseWarehouseSection(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Simulation
		hasError       bool
		errorKind      ParserErrorKind
	}

	var testCases = []testCase{
		{
			input:     []string{""},
			hasError:  true,
			errorKind: invalidNumberOfTokens,
		},
		{
			input:     []string{"koi", "22", "45"},
			hasError:  true,
			errorKind: invalidUnsignedInteger,
		},
		{
			input:     []string{"33", "pheur", "34"},
			hasError:  true,
			errorKind: invalidUnsignedInteger,
		},
		{
			input:     []string{"33", "433", "!"},
			hasError:  true,
			errorKind: invalidUnsignedInteger,
		},
		{
			input:    []string{"453", "4952", "34"},
			hasError: false,
			expectedOutput: Simulation{
				cycle: 34,
				warehouse: Warehouse{
					width:  453,
					length: 4952,
				},
			},
		},
	}

	for _, testCase := range testCases {
		var simulation, err = parseWarehouseSection(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, simulation, testCase.expectedOutput)
		}
	}
}

func TestParseParcel(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Parcel
		hasError       bool
		errorKind      ParserErrorKind
	}

	var testCases = []testCase{
		{
			input:     []string{""},
			hasError:  true,
			errorKind: invalidNumberOfTokens,
		},
		{
			input:     []string{"", "23", "30", "yellow"},
			hasError:  true,
			errorKind: invalidTokenLength,
		},
		{
			input:     []string{"parcel", "x_coord", "30", "yellow"},
			hasError:  true,
			errorKind: invalidUnsignedInteger,
		},
		{
			input:     []string{"parcel", "1", "y_coord", "yellow"},
			hasError:  true,
			errorKind: invalidUnsignedInteger,
		},
		{
			input:     []string{"parcel", "1", "1", ""},
			hasError:  true,
			errorKind: invalidWeight,
		},
		{
			input: []string{"parcel", "1", "1", "yellow"},
			expectedOutput: Parcel{
				name: "parcel",
				Coordinate: Coordinate{
					X: 1,
					Y: 1,
				},
				weight: yellow,
			},
		},
	}

	for _, testCase := range testCases {
		var parcel, err = parseParcel(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, parcel, testCase.expectedOutput)
		}
	}
}

func TestParseForklift(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Forklift
		hasError       bool
		errorKind      ParserErrorKind
	}

	var testCases = []testCase{
		{
			input: []string{"forklift", "2", "3"},
			expectedOutput: Forklift{
				name: "forklift",
				Coordinate: Coordinate{
					X: 2,
					Y: 3,
				},
			},
		},
	}

	for _, testCase := range testCases {
		var forklift, err = parseForklift(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, forklift, testCase.expectedOutput)
		}
	}
}

func TestParseTruck(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Truck
		hasError       bool
		errorKind      ParserErrorKind
	}

	var testCases = []testCase{
		{
			input: []string{"truck", "2", "3", "4000", "5"},
			expectedOutput: Truck{
				name: "truck",
				Coordinate: Coordinate{
					X: 2,
					Y: 3,
				},
				max_weight: 4000,
				available:  5,
			},
		},
	}

	for _, testCase := range testCases {
		var truck, err = parseTruck(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, truck, testCase.expectedOutput)
		}
	}
}

func TestParseWeight(t *testing.T) {
	type testCase struct {
		input          string
		expectedOutput Weight
		hasError       bool
	}

	var testCases = []testCase{
		{
			input:          "yellow",
			expectedOutput: yellow,
		},
		{
			input:          "green",
			expectedOutput: green,
		},
		{
			input:          "blue",
			expectedOutput: blue,
		},
		{
			input:    "aronge",
			hasError: true,
		},
	}

	for _, testCase := range testCases {
		var nb, err = parseWeight(testCase.input)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, nb, testCase.expectedOutput)
		}
	}
}

func TestParseNumericField(t *testing.T) {
	type testCase struct {
		input          string
		expectedOutput uint32
		hasError       bool
		fieldName      string
	}

	var testCases = []testCase{
		{
			input:          "2355",
			expectedOutput: 2355,
		},
		{
			input:     "-234",
			hasError:  true,
			fieldName: "fieldName",
		},
	}

	for _, testCase := range testCases {
		var nb, err = parseUint32Field(testCase.input)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, nb, testCase.expectedOutput)
		}
	}
}
