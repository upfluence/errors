package errors

import "github.com/upfluence/errors/tags"

func WithTags(err error, vs map[string]interface{}) error {
	return tags.WithTags(err, vs)
}
