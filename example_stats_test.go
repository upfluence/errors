package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWithStatus() {
	// Create an error and attach a status
	err := errors.New("request processing failed")
	errWithStatus := errors.WithStatus(err, "failed")
	
	fmt.Println(errWithStatus)
	// Output: request processing failed
}
