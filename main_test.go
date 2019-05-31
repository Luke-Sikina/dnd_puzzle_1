package main

import (
	"math"
	"testing"
)

func TestGenerateOperators0(t *testing.T) {
	testCases := []int{0, 1, 100}
	for index, input := range testCases {
		var operators, err = GenerateOperators(input)

		if err != nil {
			t.Errorf("Test case %d: GenerateOperators(%d) expected no error, got %s", index, input, err)
		}

		if len(operators) != input {
			t.Errorf(
				"Test case %d: GenerateOperators(%d) expected a slice of length %d, got one of length %d",
				index, input, input, len(operators))
		}
	}
}

func TestGenerateOperatorsBadInput(t *testing.T) {
	testCases := []int{-1, math.MinInt64, 101, math.MaxInt64}

	for index, input := range testCases {
		var operators, err = GenerateOperators(input)

		if err == nil {
			t.Errorf("Test case %d: GenerateOperators(%d) expected error, got nil", index, input)
		}

		if operators != nil {
			t.Errorf(
				"Test case %d: GenerateOperators(%d) expected a nil slice, got non nil slice %v",
				index, input, operators)
		}
	}
}
