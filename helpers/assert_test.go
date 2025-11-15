package helpers_test

import (
	"testing"

	"github.com/joseluisq/enve/helpers"
	"github.com/stretchr/testify/assert"
)

type mockTest struct {
	*testing.T
	failed bool
}

func (t *mockTest) Errorf(format string, args ...interface{}) {
	t.failed = true
}

func (t *mockTest) FailNow() {
	t.failed = true
}

func (t *mockTest) Fail() {
	t.failed = true
}

func TestElementsContain(t *testing.T) {
	tests := []struct {
		name     string
		t        *mockTest
		listA    any
		listB    any
		expected bool
	}{
		{
			name:     "should return true when listB is a subset of listA",
			listA:    []int{1, 2, 3, 4},
			listB:    []int{2, 4},
			expected: true,
		},
		{
			name:     "should return true when lists are identical",
			listA:    []string{"a", "b", "c"},
			listB:    []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "should return true when listB is empty",
			listA:    []int{1, 2, 3},
			listB:    []int{},
			expected: true,
		},
		{
			name:     "should return true when both lists are empty",
			listA:    []any{},
			listB:    []any{},
			expected: true,
		},
		{
			name:     "should return true with duplicate elements",
			listA:    []int{1, 2, 2, 3},
			listB:    []int{2, 2},
			expected: true,
		},
		{
			name:  "should return false when listB has elements not in listA",
			listA: []int{1, 2, 3},
			listB: []int{2, 5},
		},
		{
			name:  "should return false when listB requires more duplicates than in listA",
			listA: []int{1, 2, 3},
			listB: []int{2, 2},
		},
		{
			name:  "should return false when listA is not a slice",
			listA: "not a slice",
			listB: []int{1},
		},
		{
			name:  "should return false when listB is not a slice",
			listA: []int{1},
			listB: map[int]int{1: 1},
		},
		{
			name:  "should return false when listA is empty and listB is not",
			listB: []int{1},
		},
		{
			name:  "should return false when inputs are not of type slice",
			listA: "something",
			listB: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTest := &mockTest{T: t}
			actual := helpers.ElementsContain(mockTest, tt.listA, tt.listB)

			assert.Equal(t, tt.expected, actual, "ElementsContain should return the expected result")

			// Check if the test failed when it was expected to
			if !tt.expected && !mockTest.failed {
				assert.Fail(t, "Expected a test failure, but it passed.")
			}
		})
	}
}
