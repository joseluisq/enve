package env

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Text(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice
		expected string
	}{
		{
			name: "should return an empty string for an empty slice",
		},
		{
			name:     "should return a single line for a single element slice",
			input:    Slice{"KEY=value"},
			expected: "KEY=value",
		},
		{
			name:     "should join multiple elements with newlines",
			input:    Slice{"KEY1=value1", "KEY2=value2", "KEY3=value3"},
			expected: "KEY1=value1\nKEY2=value2\nKEY3=value3",
		},
		{
			name:     "should handle elements with no value",
			input:    Slice{"KEY_ONLY", "KEY=value"},
			expected: "KEY_ONLY\nKEY=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Text(), "Text output should match")
		})
	}
}

func TestSlice_Environ(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice
		expected Environment
	}{
		{
			name: "should return an empty Environment for an empty slice",
		},
		{
			name:  "should correctly parse valid key-value pairs",
			input: Slice{"KEY1=value1", "KEY2=value2"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "KEY1", Value: "value1"},
					{Name: "KEY2", Value: "value2"},
				},
			},
		},
		{
			name:  "should ignore invalid pairs (no equals sign)",
			input: Slice{"INVALID_KEY", "KEY=value"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "KEY", Value: "value"},
				},
			},
		},
		{
			name:  "should handle empty values",
			input: Slice{"EMPTY_VALUE=", "KEY=value"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "EMPTY_VALUE", Value: ""},
					{Name: "KEY", Value: "value"},
				},
			},
		},
		{
			name:  "should handle values with equals signs",
			input: Slice{"URL=http://example.com?param=value", "KEY=value"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "URL", Value: "http://example.com?param=value"},
					{Name: "KEY", Value: "value"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input.Environ()
			assert.Equal(t, tt.expected, actual, "Environ output should match")
		})
	}
}

func TestSlice_JSON(t *testing.T) {
	tests := []struct {
		name        string
		input       Slice
		expected    Environment
		expectedErr error
	}{
		{
			name: "should return empty JSON array for an empty slice",
		},
		{
			name:  "should return correct JSON for valid key-value pairs",
			input: Slice{"KEY1=value1", "KEY2=value2"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "KEY1", Value: "value1"},
					{Name: "KEY2", Value: "value2"},
				},
			},
		},
		{
			name:  "should ignore invalid pairs in JSON output",
			input: Slice{"INVALID_KEY", "VALID=value"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "VALID", Value: "value"},
				},
			},
		},
		{
			name:  "should return an error when parsing invalid JSON",
			input: Slice{"null"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualJSON, err := tt.input.JSON()

			if tt.expectedErr != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.Equal(t, tt.expectedErr.Error(), err.Error(), "Error message mismatch")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")

				var actual Environment
				if err := json.Unmarshal(actualJSON, &actual); tt.expectedErr != nil {
					assert.Error(t, err, "Expected an error but got none")
					assert.Equal(t, err.Error(), tt.expectedErr.Error(), "Error message mismatch")
				} else {
					assert.Equal(t, tt.expected, actual, "JSON output should match")
				}
			}
		})
	}
}

func TestSlice_XML(t *testing.T) {
	tests := []struct {
		name        string
		input       Slice
		expected    Environment
		expectedErr error
	}{
		{
			name: "should return empty XML array for an empty slice",
		},
		{
			name:  "should return correct XML for valid key-value pairs",
			input: Slice{"KEY1=value1", "KEY2=value2"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "KEY1", Value: "value1"},
					{Name: "KEY2", Value: "value2"},
				},
			},
		},
		{
			name:  "should ignore invalid pairs in XML output",
			input: Slice{"INVALID_KEY", "VALID=value"},
			expected: Environment{
				Env: []EnvironmentVar{
					{Name: "VALID", Value: "value"},
				},
			},
		},
		{
			name:  "should return an error when parsing invalid XML",
			input: Slice{"null"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualXML, err := tt.input.XML()

			if tt.expectedErr != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.Equal(t, tt.expectedErr.Error(), err.Error(), "Error message mismatch")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")

				var actual Environment
				if err := xml.Unmarshal(actualXML, &actual); tt.expectedErr != nil {
					assert.Error(t, err, "Expected an error but got none")
					assert.Equal(t, err.Error(), tt.expectedErr.Error(), "Error message mismatch")
				} else {
					assert.Equal(t, tt.expected, actual, "XML output should match")
				}
			}
		})
	}
}
