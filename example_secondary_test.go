package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWithSecondaryError() {
	// Simulate a primary error
	dbErr := errors.New("failed to save data to database")

	// Simulate a secondary error that occurred while handling the primary error
	cacheErr := errors.New("failed to invalidate cache")

	// Attach the secondary error for additional context
	err := errors.WithSecondaryError(dbErr, cacheErr)

	fmt.Println(err)
	// Output: failed to save data to database [ with secondary error: failed to invalidate cache]
}
