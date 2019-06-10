package main

import (
	"fmt"
	"testing"
)

func CreateTestGenerator(data []interface{}) chan interface{} {
	ch := make(chan interface{}, 100) //unbuffered interfaces block on <-

	go func() {
		defer close(ch)
		for _, element := range data {
			ch <- element
		}
	}()

	return ch
}

type ToSliceCase struct {
	Data     []interface{}
	Expected []interface{}
}

func TestFilterable_ToSlice(t *testing.T) {
	testCases := []ToSliceCase{
		{
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 3},
		}, {
			[]interface{}{},
			[]interface{}{},
		}, {
			[]interface{}{'a', 'b', 'c'},
			[]interface{}{'a', 'b', 'c'},
		},
	}

	for index, testCase := range testCases {
		generator := CreateTestGenerator(testCase.Data)
		filterable := NewFilterable(generator)

		actual := filterable.ToSlice()

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", testCase.Expected) {
			t.Errorf("Test case %d: ToSlice(%v) expected: %v, got: %v",
				index, testCase.Data, testCase.Expected, actual)
		}
	}
}

func AcceptEvensPredicate(toExamine interface{}) bool {
	asInt := toExamine.(int)
	return asInt%2 == 0
}

func AcceptEverythingPredicate(_ interface{}) bool {
	return true
}

func AcceptNothingPredicate(_ interface{}) bool {
	return false
}

type FilterCase struct {
	FirstPredicate  Predicate
	SecondPredicate Predicate
	Data            []interface{}
	Expected        []interface{}
}

func TestFilterable_Filter(t *testing.T) {
	testCases := []FilterCase{
		{
			AcceptEverythingPredicate,
			AcceptEverythingPredicate,
			[]interface{}{},
			[]interface{}{},
		}, {
			AcceptEverythingPredicate,
			AcceptEverythingPredicate,
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 3},
		}, {
			AcceptEverythingPredicate,
			AcceptNothingPredicate,
			[]interface{}{1, 2, 3},
			[]interface{}{},
		}, {
			AcceptEvensPredicate,
			AcceptEverythingPredicate,
			[]interface{}{1, 2, 3},
			[]interface{}{2},
		}, {
			AcceptEvensPredicate,
			AcceptEvensPredicate,
			[]interface{}{1, 2, 3},
			[]interface{}{2},
		},
	}

	for index, testCase := range testCases {
		generator := CreateTestGenerator(testCase.Data)
		filterable := NewFilterable(generator)

		actual := filterable.
			Filter(testCase.FirstPredicate).
			Filter(testCase.SecondPredicate).
			ToSlice()

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", testCase.Expected) {
			t.Errorf("Test case %d: ToSlice(%v) expected: %v, got: %v",
				index, testCase.Data, testCase.Expected, actual)
		}
	}
}
