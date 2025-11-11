// Package helpers provide utilities for the project including testing.
package helpers

import (
	"fmt"
	"reflect"

	"github.com/stretchr/testify/assert"
)

// ElementsContain asserts that all elements in listB are contained in listA.
func ElementsContain(t assert.TestingT, listA any, listB any, msgAndArgs ...any) (ok bool) {
	aVal := reflect.ValueOf(listA)
	bVal := reflect.ValueOf(listB)

	if aVal.Kind() != reflect.Slice || bVal.Kind() != reflect.Slice {
		return assert.Fail(t, "ElementsContain only accepts slice arguments", msgAndArgs...)
	}

	// Build multiset for listA
	counts := make(map[any]int)
	for i := 0; i < aVal.Len(); i++ {
		val := aVal.Index(i).Interface()
		counts[val]++
	}

	// Check that each element in listB is present in listA
	for i := 0; i < bVal.Len(); i++ {
		val := bVal.Index(i).Interface()
		if counts[val] == 0 {
			return assert.Fail(
				t, fmt.Sprintf("Expected element %+v not found in listA: %+v", val, listA), msgAndArgs...,
			)
		}
		counts[val]--
	}

	return true
}
