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
