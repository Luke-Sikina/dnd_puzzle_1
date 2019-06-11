package main

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestGenerateOperators0(t *testing.T) {
	testCases := []int{0, 1, 15}
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
	testCases := []int{-1, math.MinInt64, 10000, math.MaxInt64}

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
		{Candidate{}, []int{2}, 2, true},
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

type FilterOperatorsCase struct {
	Unfiltered []Candidate
	Terms      []int
	Goal       int
	Expected   []Candidate
}

func TestFilterOperators(t *testing.T) {
	testCases := []FilterOperatorsCase{
		{
			[]Candidate{{PLUS}, {MINUS}, {DIVIDE}, {TIMES}},
			[]int{1, 1},
			2,
			[]Candidate{{PLUS}},
		}, {
			[]Candidate{{PLUS}, {MINUS}, {DIVIDE}, {TIMES}},
			[]int{1, 1},
			1,
			[]Candidate{{DIVIDE}, {TIMES}},
		}, {[]Candidate{{PLUS}, {MINUS}, {DIVIDE}, {TIMES}},
			[]int{1, 1},
			3,
			[]Candidate{},
		}, {
			[]Candidate{},
			[]int{1, 1},
			3,
			[]Candidate{},
		}, {
			[]Candidate{{PLUS}},
			[]int{1},
			3,
			[]Candidate{},
		}, {
			[]Candidate{{PLUS}},
			[]int{1, 1, 1},
			3,
			[]Candidate{},
		},
	}

	for index, testCase := range testCases {
		actual := FilterOperators(testCase.Unfiltered, testCase.Terms, testCase.Goal)

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", testCase.Expected) {
			t.Errorf("Test case %d: FilterOperators(%v %v %d) expected: %v, got: %v",
				index, testCase.Unfiltered, testCase.Terms, testCase.Goal, testCase.Expected, actual)
		}
	}
}

type GenerateAllOperatorsCase struct {
	OperatorCount uint
	Candidates    []Candidate
	Error         error
}

func TestOperatorGenerator(t *testing.T) {
	testCases := []GenerateAllOperatorsCase{
		{
			1,
			[]Candidate{{PLUS}, {MINUS}, {TIMES}, {DIVIDE}},
			nil,
		},
		{
			2,
			[]Candidate{
				{PLUS, PLUS}, {MINUS, PLUS}, {TIMES, PLUS}, {DIVIDE, PLUS},
				{PLUS, MINUS}, {MINUS, MINUS}, {TIMES, MINUS}, {DIVIDE, MINUS},
				{PLUS, TIMES}, {MINUS, TIMES}, {TIMES, TIMES}, {DIVIDE, TIMES},
				{PLUS, DIVIDE}, {MINUS, DIVIDE}, {TIMES, DIVIDE}, {DIVIDE, DIVIDE},
			},
			nil,
		}, {
			0,
			nil,
			errors.New("bad operator count: 1000 > operators > 0"),
		},
	}

	for index, testCase := range testCases {
		generator, actualError := OperatorGenerator(testCase.OperatorCount)
		filterable := NewFilterable(generator)
		actual := filterable.ToSlice()
		if fmt.Sprintf("%v", actualError) != fmt.Sprintf("%v", testCase.Error) {
			t.Errorf("Test case %d: GenerateAllOperators(%d) expected error: %v, got error: %v",
				index, testCase.OperatorCount, testCase.Error, actualError)
		}

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", testCase.Candidates) {
			t.Errorf("Test case %d: GenerateAllOperators(%d) expected: %v, got: %v",
				index, testCase.OperatorCount, testCase.Candidates, actual)
		}
	}
}

type GenerateTermsCase struct {
	Count    int
	Min      int
	Max      int
	Expected Terms
}

func TestGenerateTerms(t *testing.T) {
	testCases := []GenerateTermsCase{
		{1, 1, 2, Terms{[]int{1}, 0}},
		{1, 0, 1, Terms{[]int{0}, 0}},
		{3, 2, 3, Terms{[]int{2, 2, 2}, 0}},
	}

	//trying to remove as much logic from this test as possible, so the test isnt going to verify the ranges
	for index, testCase := range testCases {
		actual := GenerateTerms(testCase.Count, testCase.Min, testCase.Max)

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", testCase.Expected) {
			t.Errorf("Test case %d: GenerateTerms(%d, %d, %d) expected: %v, got %v",
				index, testCase.Count, testCase.Min, testCase.Max, testCase.Expected, actual)
		}
	}
}

func candidateGenerator(candidates []Candidate) chan interface{} {
	ch := make(chan interface{}, 5)

	go func() {
		defer close(ch)
		for _, candidate := range candidates {
			ch <- candidate
		}
	}()

	return ch
}

// Integration tests
func TestGenerateTermsUntilSingleCandidate(t *testing.T) {
	candidates := []Candidate{
		[]Operator{PLUS},
		[]Operator{MINUS},
		[]Operator{TIMES},
		[]Operator{DIVIDE},
	}
	termCount := 2
	min := 3
	max := 4
	operators := []Operator{TIMES}

	actual := GenerateTermsUntilSingleCandidate(candidateGenerator(candidates), termCount, min, max, 1, operators)

	expected := []Terms{{[]int{3, 3}, 9}}

	if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
		t.Errorf("GenerateTermsUntilSingleCandidate(%v, %d, %d, %d, %v) expected: %v, got: %v",
			candidates, termCount, min, max, operators, expected, actual)
	}
}
