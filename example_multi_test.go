package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWrapErrors() {
	var errs []error

	// Simulate validation errors
	errs = append(errs, errors.New("name is required"))
	errs = append(errs, errors.New("email is invalid"))

	if len(errs) > 0 {
		err := errors.WrapErrors(errs)
		fmt.Println(err)
	}
	// Output: [name is required, email is invalid]
}

func ExampleCombine() {
	err := errors.Combine(
		errors.New("validation failed: name"),
		errors.New("validation failed: email"),
		errors.New("validation failed: age"),
	)
	fmt.Println(err)
	// Output: [validation failed: name, validation failed: email, validation failed: age]
}

func ExampleJoin() {
	err := errors.Join(
		errors.New("first error"),
		errors.New("second error"),
	)
	fmt.Println(err)
	// Output: [first error, second error]
}
