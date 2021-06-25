package main

import (
	"fmt"
	"log"
)

type ErrorAccumulator struct {
	errors []error
}

func (a *ErrorAccumulator) accumulate(err error) {
	a.errors = append(a.errors, err)
}

func (a *ErrorAccumulator) step(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (a *ErrorAccumulator) compileOR() error {
	finalString := ""

	for _, err := range a.errors {
		finalString += fmt.Sprintln("|", err)
	}

	return fmt.Errorf("Syntax error at %s:%d:%d, expected either one of : \n%s", filePath, 0, 0, finalString)
}

func (a *ErrorAccumulator) compileFirst() error {
	return fmt.Errorf("Syntax error at %s:%d:%d, expected \n%s", filePath, 0, 0, a.errors[0])
}

func errorAcc() ErrorAccumulator {
	return ErrorAccumulator{[]error{}}
}
