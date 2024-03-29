package main

import (
	"errors"
	"flag"
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
	if terms < 0 || terms > 1000 { // these bounds make more sense in GenerateAllOperators, copied here for consistency
		return nil, errors.New("1000 > terms < 0")
	}
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	operators := make([]Operator, terms, terms)
	for i := 0; i < terms; i++ {
		operators[i] = Operator(randGenerator.Intn(4))
	}

	return operators, nil
}

func OperatorGenerator(operators uint) (chan interface{}, error) {
	ch := make(chan interface{}, 1000000) //unbuffered interfaces block on <-
	if operators <= 0 || operators > 1000 {
		close(ch) //returning the actual empty, closed channel made testing a bit more fluid
		return ch, errors.New("bad operator count: 1000 > operators > 0")
	}

	go func() {
		defer close(ch)
		var max uint = 1 << (2 * operators)
		for i := uint(0); i < max; i++ {
			candidate := make(Candidate, operators, operators)
			bits := i
			for operatorIndex := uint(0); operatorIndex < operators; operatorIndex++ {
				candidate[operatorIndex] = Operator(bits % 4)
				bits = bits >> 2
			}
			ch <- candidate
		}
	}()

	return ch, nil
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

func parseParams() (terms, min, max, attempts int, err error) {
	flag.IntVar(&terms, "terms", 3, "Number of terms (numbers) in the generated problems")
	flag.IntVar(&min, "min", 2, "Minimum value for a term (number) in the generated problems")
	flag.IntVar(&max, "max", 20, "Maximum value for a term (number) in the generated problems")
	flag.IntVar(&attempts, "att", 20, "Number of attempts (number)")
	flag.Parse()
	if min >= max {
		err = errors.New("min must be  <= max")
	}
	log.Printf("Creating a %d term problem. Each term has a min of %d and a max of %d", terms, min, max)
	return
}

func ifErrThenExit(message string, err error) {
	if err != nil {
		log.Printf(message, err)
		os.Exit(1)
	}
}

func main() {
	termCount, min, max, attempts, err := parseParams()
	ifErrThenExit("Error parsing params: %v", err)

	operators, err := GenerateOperators(termCount - 1)
	ifErrThenExit("Error generating operators: %v", err)

	generator, err := OperatorGenerator(uint(termCount - 1))
	ifErrThenExit("error generating candidates: %v", err)

	allTerms := GenerateTermsUntilSingleCandidate(generator, termCount, min, max, attempts, operators)
	fmt.Printf("For the operators: %v, the terms generated are:\n", operators)
	for _, row := range allTerms {
		fmt.Printf("%v\n", row)
	}
}

func GenerateTermsUntilSingleCandidate(operatorGen chan interface{}, termCount, min, max, maxAttempts int, actualOperators []Operator) (allTerms []Terms) {
	filterable := NewFilterable(operatorGen)
	for attempts := 0; attempts < maxAttempts; attempts++ {
		terms := GenerateTerms(termCount, min, max)
		goal := terms.Terms[0]
		for index, operator := range actualOperators {
			goal = EvaluateOperator(goal, terms.Terms[index+1], operator)
		}
		terms.Solution = goal
		allTerms = append(allTerms, terms)
		log.Printf("Created terms: %v.", terms)
		predicate := func(candidate interface{}) bool { return EvaluateCandidate(candidate.(Candidate), terms.Terms, goal) }
		filterable = filterable.Filter(predicate)
	}
	return
}
