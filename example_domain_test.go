package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWithDomain() {
	// Create an error and attach a domain for categorization
	err := errors.New("operation failed")
	errWithDomain := errors.WithDomain(err, "database")
	
	fmt.Println(errWithDomain)
	// Output: operation failed
}
