package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				Coordonate: Coordonate{
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
				Coordonate: Coordonate{
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
				Coordonate: Coordonate{
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
			assert.NotEqual(t, err, nil)
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
			assert.NotEqual(t, err, nil)
		} else {
			assert.Equal(t, nb, testCase.expectedOutput)
		}
	}
}
