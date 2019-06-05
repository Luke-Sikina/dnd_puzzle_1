package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Operator int
type Candidate []Operator
type Terms struct {
	Terms    []int
	Solution int
}

const (
	PLUS Operator = iota
	MINUS
	TIMES
	DIVIDE
)

func (terms Terms) Format(f fmt.State, c rune) {
	toPrint := fmt.Sprintf("%v = %d", terms.Terms, terms.Solution)
	_, err := f.Write([]byte(toPrint))
	if err != nil {
		log.Printf("Error formatting terms: %v", err)
	}
}

func (operator Operator) Format(f fmt.State, c rune) {
	var toWrite string
	switch operator {
	case PLUS:
		toWrite = "+"
	case MINUS:
		toWrite = "-"
	case TIMES:
		toWrite = "*"
	case DIVIDE:
		toWrite = "/"
	}
	_, err := f.Write([]byte(toWrite))
	if err != nil {
		log.Printf("Error formatting operator: %v", err)
	}
}

func GenerateOperators(terms int) ([]Operator, error) {
	if terms < 0 || terms > 15 { // these bounds make more sense in GenerateAllOperators, copied here for consistency
		return nil, errors.New("terms < 0")
	}
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	operators := make([]Operator, terms, terms)
	for i := 0; i < terms; i++ {
		operators[i] = Operator(randGenerator.Intn(4))
	}

	return operators, nil
}

func GenerateAllOperators(operators uint) ([]Candidate, error) {
	if operators <= 0 || operators > 15 {
		return nil, errors.New("bad operator count: 16 > operators > 0")
	}

	//there are 4 operators, so each operator is 2 bytes in an unsigned int
	var max uint = 1 << (2 * operators)
	candidates := make([]Candidate, 1<<(2*operators), 1<<(2*operators))
	for i := uint(0); i < max; i++ {
		candidate := make(Candidate, operators, operators)
		bits := i
		for operatorIndex := uint(0); operatorIndex < operators; operatorIndex++ {
			candidate[operatorIndex] = Operator(bits % 4)
			bits = bits >> 2
			candidates[i] = candidate
		}
	}

	return candidates, nil
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

func GenerateTerms(count, min, max int) Terms {
	terms := make([]int, count, count)
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		terms[i] = randGenerator.Intn(max-min) + min
	}

	return Terms{terms, 0}
}

func main() {
	termCount := 12
	operators, err := GenerateOperators(termCount - 1)

	if err != nil {
		os.Exit(1)
	}

	allOperators, err := GenerateAllOperators(uint(termCount - 1))

	if err != nil {
		os.Exit(1)
	}

	allTerms := GenerateTermsUntilSingleCandidate(allOperators, termCount, operators)
	fmt.Printf("For the operators: %v, the terms generated are:\n", operators)
	for _, row := range allTerms {
		fmt.Printf("%v\n", row)
	}
}

func GenerateTermsUntilSingleCandidate(allOperators []Candidate, termCount int, actualOperators []Operator) (allTerms []Terms) {
	for len(allOperators) > 1 {
		terms := GenerateTerms(termCount, 2, 20)
		goal := terms.Terms[0]
		for index, operator := range actualOperators {
			goal = EvaluateOperator(goal, terms.Terms[index+1], operator)
		}
		oldLength := len(allOperators)
		allOperators = FilterOperators(allOperators, terms.Terms, goal)
		if oldLength > len(allOperators) {
			terms.Solution = goal
			allTerms = append(allTerms, terms)
			log.Printf("Created terms: %v, currently %d valid candidates.", terms, len(allOperators))
		} else {
			log.Printf("Rejected terms: %v, did not reduce candidates.", terms)
		}
	}
	return
}
