package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySimulationValidity(t *testing.T) {
	type testCase struct {
		input    Simulation
		hasError bool
	}

	testCases := []testCase{
		{
			input:    Simulation{},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					length: 5,
					width:  5,
					forklifts: []forklift{
						{name: "forklift"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					width:  1,
					length: 2,
					forklifts: []forklift{
						{name: "forklift", coordinate: coordinate{X: 4}},
					},
					trucks: []truck{
						{name: "truck"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					width:  1,
					length: 2,
					forklifts: []forklift{
						{name: "forklift", coordinate: coordinate{X: 0, Y: 0}},
					},
					trucks: []truck{
						{name: "truck", coordinate: coordinate{X: 0, Y: 0}},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					width:  2,
					length: 1,
					forklifts: []forklift{
						{name: "entity", coordinate: coordinate{X: 1}},
					},
					trucks: []truck{
						{name: "entity"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					width:  6,
					length: 6,
					forklifts: []forklift{
						{name: "forklift", coordinate: coordinate{X: 1}},
					},
					trucks: []truck{
						{name: "truck", coordinate: coordinate{X: 4, Y: 4}},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				warehouse: warehouse{
					width:  2,
					length: 1,
					forklifts: []forklift{
						{name: "forklift", coordinate: coordinate{X: 1}},
					},
					trucks: []truck{
						{name: "truck"},
					},
				},
			},
			hasError: false,
		},
	}

	for _, testCase := range testCases {
		err := VerifySimulationValidity(testCase.input)

		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.False(t, testCase.hasError)
		}
	}
}
