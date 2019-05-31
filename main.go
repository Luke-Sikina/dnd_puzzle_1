package main

import (
	"errors"
	"math/rand"
	"time"
)

type Operator int

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

func main() {

}
