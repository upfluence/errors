//go:build go1.18

package errors_test

import (
	"fmt"

	"github.com/upfluence/errors"
)

type MyError struct {
	message string
}

func (e MyError) Error() string { return e.message }

type MyErrorWithCode struct {
	Code    int
	Message string
}

func (e MyErrorWithCode) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

func ExampleIsOfType() {
	err := errors.Wrap(MyError{message: "my error"}, "wrapped error")

	if errors.IsOfType[MyError](err) {
		fmt.Println("MyError was found in the error chain")
	}
	// Output: MyError was found in the error chain
}

func ExampleAsType() {
	err := errors.Wrap(MyErrorWithCode{Code: 404, Message: "not found"}, "request failed")

	if myErr, ok := errors.AsType[MyErrorWithCode](err); ok {
		fmt.Println(myErr.Code)
	}
	// Output: 404
}
