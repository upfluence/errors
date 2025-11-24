package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

func ExampleWithTags() {
	userID := "user123"
	
	// Create an error and attach tags for additional context
	err := errors.New("failed to fetch user")
	errWithTags := errors.WithTags(err, map[string]interface{}{
		"user_id":   userID,
		"operation": "fetch",
	})
	
	fmt.Println(errWithTags)
	// Output: failed to fetch user
}
