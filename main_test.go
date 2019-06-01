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
		operators, err := GenerateOperators(input)

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

type EvaluateCandidateCase struct {
	Candidate Candidate
	Terms     []int
	Goal      int
	Expected  bool
}

func TestEvaluateCandidate(t *testing.T) {
	testCases := []EvaluateCandidateCase{
		{Candidate{PLUS}, []int{2, 1}, 2, false},
		{Candidate{TIMES}, []int{2, 1}, 2, true},
		{Candidate{MINUS}, []int{2, 1}, 2, false},
		{Candidate{DIVIDE}, []int{2, 1}, 2, true},
		{Candidate{PLUS, TIMES}, []int{2, 1, 4}, 12, true},
		{Candidate{PLUS, TIMES}, []int{2, 1, 4}, 11, false},
		{Candidate{PLUS}, []int{2, 1, 4}, 12, false},
		{Candidate{}, []int{}, 12, false},
	}

	for index, testCase := range testCases {
		actual := EvaluateCandidate(testCase.Candidate, testCase.Terms, testCase.Goal)

		if actual != testCase.Expected {
			t.Errorf("Test case %d: EvaluateCandidate(%v, %v, %d) expected %v got %v",
				index, testCase.Candidate, testCase.Terms, testCase.Goal, testCase.Expected, actual)
		}
	}
}

type EvaluateOperatorCase struct {
	Operator Operator
	First    int
	Second   int
	Expected int
}

func TestEvaluateOperator(t *testing.T) {
	testCases := []EvaluateOperatorCase{
		{PLUS, 1, 1, 2},
		{MINUS, 1, 1, 0},
		{TIMES, 1, 1, 1},
		{DIVIDE, 1, 1, 1},
		{PLUS, 10, 2, 12},
		{MINUS, 10, 2, 8},
		{TIMES, 10, 2, 20},
		{DIVIDE, 10, 2, 5},
		{DIVIDE, 1, 0, math.MaxInt64},
	}

	for index, testCase := range testCases {
		actual := EvaluateOperator(testCase.First, testCase.Second, testCase.Operator)

		if actual != testCase.Expected {
			t.Errorf("Test case %d: EvaluateOperator(%d %d %d) expected: %d, got: %d",
				index, testCase.First, testCase.Second, testCase.Operator, testCase.Expected, actual)
		}
	}
}
