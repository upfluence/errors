package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleOpaque() {
	// Create an internal error that we want to hide
	internalErr := errors.New("database connection pool exhausted")
	
	// Make it opaque to prevent callers from inspecting internal details
	err := errors.Opaque(internalErr)
	
	// The error message is still available
	fmt.Println(err)
	
	// But unwrapping won't work - it returns nil
	fmt.Println(errors.Unwrap(err) == nil)
	// Output: database connection pool exhausted
	// true
}
