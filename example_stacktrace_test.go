package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWithStack() {
	// Simulate an error from an external library
	externalErr := fmt.Errorf("external library error")

	// Add stack trace at current location
	err := errors.WithStack(externalErr)

	fmt.Println(err)
	// Output: external library error
}

// Result represents a simple computation result
type Result struct {
	Value int
}

// mockCompute simulates an external library function
func mockCompute(input int) (Result, error) {
	if input < 0 {
		return Result{}, fmt.Errorf("negative input not allowed")
	}
	return Result{Value: input * 2}, nil
}

func ExampleWithStack2() {
	// Demonstrates using WithStack2 to add stack traces inline while returning values
	compute := func(input int) (Result, error) {
		// WithStack2 adds stack trace in a single line
		return errors.WithStack2(mockCompute(input))
	}

	// Success case - error is nil
	result, err := compute(5)
	fmt.Printf("Result: %+v, Error: %v\n", result, err)

	// Error case - error gets stack trace added
	result, err = compute(-1)
	fmt.Printf("Result: %+v, Error: %v\n", result, err)

	// Output:
	// Result: {Value:10}, Error: <nil>
	// Result: {Value:0}, Error: negative input not allowed
}

func ExampleWithFrame() {
	// This helper function wraps errors with the caller's location
	wrapError := func(err error) error {
		// Skip 1 frame to capture the caller's location instead of this function
		return errors.WithFrame(err, 1)
	}

	err := wrapError(fmt.Errorf("original error"))
	fmt.Println(err)
	// Output: original error
}
