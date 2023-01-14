package parsing

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
				Warehouse: Warehouse{
					Length: 5,
					Width:  5,
					Forklifts: []Forklift{
						{Name: "forklift"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				Warehouse: Warehouse{
					Width:  1,
					Length: 2,
					Forklifts: []Forklift{
						{Name: "forklift", coordinate: coordinate{X: 4}},
					},
					Trucks: []Truck{
						{Name: "truck"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				Warehouse: Warehouse{
					Width:  1,
					Length: 2,
					Forklifts: []Forklift{
						{Name: "forklift", coordinate: coordinate{X: 0, Y: 0}},
					},
					Trucks: []Truck{
						{Name: "truck", coordinate: coordinate{X: 0, Y: 0}},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				Warehouse: Warehouse{
					Width:  2,
					Length: 1,
					Forklifts: []Forklift{
						{Name: "entity", coordinate: coordinate{X: 1}},
					},
					Trucks: []Truck{
						{Name: "entity"},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				Warehouse: Warehouse{
					Width:  6,
					Length: 6,
					Forklifts: []Forklift{
						{Name: "forklift", coordinate: coordinate{X: 1}},
					},
					Trucks: []Truck{
						{Name: "truck", coordinate: coordinate{X: 4, Y: 4}},
					},
				},
			},
			hasError: true,
		},
		{
			input: Simulation{
				Warehouse: Warehouse{
					Width:  2,
					Length: 1,
					Forklifts: []Forklift{
						{Name: "forklift", coordinate: coordinate{X: 1}},
					},
					Trucks: []Truck{
						{Name: "truck"},
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
