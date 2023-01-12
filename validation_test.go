package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySimulationValidity(t *testing.T) {
	type testCase struct {
		input    Simulation
		hasError bool
	}

	var testCases = []testCase{
		{
			input:    Simulation{},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: Warehouse{
					forklifts: []Forklift{
						{name: "forklift"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: Warehouse{
					width:  1,
					length: 1,
					forklifts: []Forklift{
						{name: "forklift", Coordinate: Coordinate{X: 4}},
					},
					trucks: []Truck{
						{name: "truck"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: Warehouse{
					width:  1,
					length: 1,
					forklifts: []Forklift{
						{name: "forklift", Coordinate: Coordinate{X: 0, Y: 0}},
					},
					trucks: []Truck{
						{name: "truck", Coordinate: Coordinate{X: 0, Y: 0}},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: Warehouse{
					width:  2,
					length: 1,
					forklifts: []Forklift{
						{name: "entity", Coordinate: Coordinate{X: 1}},
					},
					trucks: []Truck{
						{name: "entity"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: Warehouse{
					width:  2,
					length: 1,
					forklifts: []Forklift{
						{name: "forklift", Coordinate: Coordinate{X: 1}},
					},
					trucks: []Truck{
						{name: "truck"},
					},
				},
			},
			hasError: false,
		},
	}

	for _, testCase := range testCases {
		var err = verifySimulationValidity(testCase.input)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.False(t, testCase.hasError)
		}
	}
}
