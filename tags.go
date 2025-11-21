package errors

import "github.com/upfluence/errors/tags"

// WithTags attaches key-value tags to the error for additional context and adds a stack frame.
func WithTags(err error, vs map[string]interface{}) error {
	return WithFrame(tags.WithTags(err, vs), 1)
}
