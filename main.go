package main

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Operator int
type Candidate []Operator

const (
	PLUS Operator = iota
	MINUS
	TIMES
	DIVIDE
)

func GenerateOperators(terms int) ([]Operator, error) {
	if terms < 0 || terms > 100 { // sanity check, more than 100 terms seems like way too much
		return nil, errors.New("terms < 0")
	}
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	operators := make([]Operator, terms, terms)
	for i := 0; i < terms; i++ {
		operators[i] = Operator(randGenerator.Intn(4))
	}

	return operators, nil
}

func FilterOperators(unfiltered []Candidate, terms []int, goal int) (filtered []Candidate) {
	for _, candidate := range unfiltered {
		if EvaluateCandidate(candidate, terms, goal) {
			filtered = append(filtered, candidate)
		}
	}
	return
}

func EvaluateCandidate(candidate Candidate, terms []int, goal int) bool {
	if len(candidate)+1 != len(terms) {
		return false
	}

	if len(terms) == 1 {
		return terms[0] == goal
	}

	total := terms[0]

	for index, operator := range candidate {
		total = EvaluateOperator(total, terms[index+1], operator)
	}
	return total == goal
}

func EvaluateOperator(first int, second int, operator Operator) int {
	switch operator {
	case PLUS:
		return first + second
	case MINUS:
		return first - second
	case TIMES:
		return first * second
	case DIVIDE:
		if second == 0 {
			return math.MaxInt64
		}
		return first / second
	default:
		return first // this should be unreachable, but the compiler requires it
	}
}

func main() {

}
