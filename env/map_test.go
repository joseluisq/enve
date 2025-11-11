package env_test

import (
	"testing"

	"github.com/joseluisq/enve/env"
	"github.com/joseluisq/enve/helpers"
)

func TestMap_Array(t *testing.T) {
	tests := []struct {
		name     string
		input    env.Map
		expected []string
	}{
		{
			name:     "should return an empty slice for an empty map",
			input:    env.Map{},
			expected: []string{},
		},
		{
			name: "should return array with valid values",
			input: env.Map{
				"KEY2": "value2",
				"KEY1": "value1",
			},
			expected: []string{
				"KEY2=value2",
				"KEY1=value1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input.Array()
			helpers.ElementsContain(t, tt.expected, actual, "Array output should match")
		})
	}
}
