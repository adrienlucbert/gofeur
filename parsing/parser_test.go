package parsing

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

	testCases := []testCase{
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
				Cycle: 1000,
				Warehouse: Warehouse{
					Width:  5,
					Length: 4,
					Parcels: []Parcel{
						{
							Name:       "colis_a_livrer",
							coordinate: coordinate{X: 2, Y: 1},
							Weight:     green,
						},
						{
							Name:       "paquet",
							coordinate: coordinate{X: 2, Y: 2},
							Weight:     blue,
						},
						{
							Name:       "deadpool",
							coordinate: coordinate{X: 0, Y: 3},
							Weight:     yellow,
						},
						{
							Name: "colère_DU_dragon",
							coordinate: coordinate{
								X: 4,
								Y: 1,
							},
							Weight: green,
						},
					},
					Forklifts: []Forklift{
						{
							Name:       "transpalette_1",
							coordinate: coordinate{X: 0, Y: 5},
						},
					},
					Trucks: []Truck{
						{
							Name:       "camion_b",
							coordinate: coordinate{X: 3, Y: 4},
							MaxWeight:  4000,
							Available:  5,
						},
						{
							Name:       "camion_a",
							coordinate: coordinate{X: 2, Y: 2},
							MaxWeight:  4007,
							Available:  4,
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
				Cycle: 243,
				Warehouse: Warehouse{
					Width: 10, Length: 50,
					Forklifts: []Forklift{
						{
							Name:       "forklift",
							coordinate: coordinate{X: 1, Y: 10},
						},
					},
					Trucks: []Truck{
						{
							Name:       "truck",
							coordinate: coordinate{X: 0, Y: 5},
							MaxWeight:  10000,
							Available:  60,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		input := strings.Join(testCase.input, "\n")
		reader := strings.NewReader(input)
		simul, err := parseFromReader(reader)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, simul, testCase.expectedOutput)
		}
	}
}

func TestParseWarehouseSection(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Simulation
		hasError       bool
		errorKind      parserErrorKind
	}

	testCases := []testCase{
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
			input:     []string{"1", "1", "9"},
			hasError:  true,
			errorKind: invalidCycleNumber,
		},
		{
			input:     []string{"1", "1", "100001"},
			hasError:  true,
			errorKind: invalidCycleNumber,
		},
		{
			input:    []string{"453", "4952", "34"},
			hasError: false,
			expectedOutput: Simulation{
				Cycle: 34,
				Warehouse: Warehouse{
					Width:  453,
					Length: 4952,
				},
			},
		},
	}

	for _, testCase := range testCases {
		simul, err := parseWarehouseSection(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, simul, testCase.expectedOutput)
		}
	}
}

func TestParseParcel(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Parcel
		hasError       bool
		errorKind      parserErrorKind
	}

	testCases := []testCase{
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
				Name: "parcel",
				coordinate: coordinate{
					X: 1,
					Y: 1,
				},
				Weight: yellow,
			},
		},
	}

	for _, testCase := range testCases {
		parcel, err := parseParcel(testCase.input)

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
		errorKind      parserErrorKind
	}

	testCases := []testCase{
		{
			input: []string{"forklift", "2", "3"},
			expectedOutput: Forklift{
				Name: "forklift",
				coordinate: coordinate{
					X: 2,
					Y: 3,
				},
			},
		},
	}

	for _, testCase := range testCases {
		flt, err := parseForklift(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, flt, testCase.expectedOutput)
		}
	}
}

func TestParseTruck(t *testing.T) {
	type testCase struct {
		input          []string
		expectedOutput Truck
		hasError       bool
		errorKind      parserErrorKind
	}

	testCases := []testCase{
		{
			input: []string{"truck", "2", "3", "4000", "5"},
			expectedOutput: Truck{
				Name: "truck",
				coordinate: coordinate{
					X: 2,
					Y: 3,
				},
				MaxWeight: 4000,
				Available: 5,
			},
		},
	}

	for _, testCase := range testCases {
		lorry, err := parseTruck(testCase.input)

		if testCase.hasError {
			assert.Equal(t, err.Kind(), testCase.errorKind)
		} else {
			assert.Equal(t, lorry, testCase.expectedOutput)
		}
	}
}

func TestParseWeight(t *testing.T) {
	type testCase struct {
		input          string
		expectedOutput weight
		hasError       bool
	}

	testCases := []testCase{
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
		nb, err := parseWeight(testCase.input)

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

	testCases := []testCase{
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
		nb, err := parseUint32Field(testCase.input)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, nb, testCase.expectedOutput)
		}
	}
}
