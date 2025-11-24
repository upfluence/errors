// Package tags provides utilities for extracting tags from errors.
//
// This package allows errors to carry arbitrary key-value metadata (tags)
// that can be used for logging, metrics, or error reporting. It traverses
// the error chain to collect all tags, with outer tags taking precedence
// over inner tags when keys conflict.
package tags

import "github.com/upfluence/errors/base"

// GetTags extracts all tags from an error by traversing the error chain.
// Returns nil if no tags are found.
// When multiple errors in the chain have the same tag key, the outermost value is used.
func GetTags(err error) map[string]interface{} {
	var tags map[string]interface{}

	for {
		if err == nil {
			break
		}

		if t, ok := err.(interface{ Tags() map[string]interface{} }); ok {
			ts := t.Tags()

			if len(ts) > 0 && tags == nil {
				tags = make(map[string]interface{}, len(ts))
			}

			for k, v := range ts {
				if _, ok := tags[k]; !ok {
					tags[k] = v
				}
			}
		}

		err = base.UnwrapOnce(err)
	}

	return tags
}
