// Package multi provides support for handling multiple errors as a single error.
//
// This package allows combining multiple errors into one, which is useful for
// operations that can fail in multiple ways or when collecting errors from
// concurrent operations. It automatically flattens nested multi-errors and
// filters out nil errors.
package multi

import (
	"errors"
	"strings"

	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/tags"
)

type multiError []error

func (errs multiError) Unwrap() []error { return errs }
func (errs multiError) Errors() []error { return errs }

func (errs multiError) Error() string {
	var b strings.Builder

	b.WriteRune('[')

	for i, err := range errs {
		b.WriteString(err.Error())

		if i < len(errs)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteRune(']')

	return b.String()
}

func (errs multiError) Tags() map[string]interface{} {
	var allTags map[string]interface{}

	for _, err := range errs {
		ts := tags.GetTags(err)

		if len(ts) > 0 && allTags == nil {
			allTags = make(map[string]interface{}, len(ts))
		}

		for k, v := range ts {
			if _, ok := allTags[k]; !ok {
				allTags[k] = v
			}
		}
	}

	return allTags
}

// Wrap combines multiple errors into a single error.
// Returns nil if all errors are nil, returns the single error if only one is non-nil,
// or returns a multiError containing all non-nil errors.
// Automatically flattens nested multiErrors.
func Wrap(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}

	var merrs []error

	for _, err := range errs {
		if err == nil {
			continue
		}

		merrs = append(merrs, ExtractErrors(err)...)
	}

	switch len(merrs) {
	case 0:
		return nil
	case 1:
		return merrs[0]
	default:
		return multiError(merrs)
	}
}

// Combine combines multiple errors into a single error.
// This is an alias for Wrap that accepts variadic arguments.
func Combine(errs ...error) error { return Wrap(errs) }

// MultiError is an interface for errors that contain multiple errors.
type MultiError interface {
	Errors() []error
}

// ExtractErrors recursively extracts all errors from an error,
// flattening any MultiError instances. Returns nil if err is nil.
func ExtractErrors(err error) []error {
	if err == nil {
		return nil
	}

	merr, ok := err.(MultiError)

	if ok {
		return merr.Errors()
	}

	nerr := base.UnwrapOnce(err)

	if nerr == nil {
		return []error{err}
	}

	errs := ExtractErrors(nerr)

	if len(errs) == 1 && errors.Is(errs[0], nerr) {
		return []error{err}
	}

	return errs
}
