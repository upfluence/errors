package multi

import (
	"strings"

	"github.com/upfluence/errors/tags"
)

type multiError []error

func (errs multiError) Errors() []error {
	return errs
}

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

func Wrap(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}

	var merrs []error

	for _, err := range errs {
		if e, ok := err.(interface{ Errors() []error }); ok {
			merrs = append(merrs, e.Errors()...)
		} else {
			merrs = append(merrs, err)
		}
	}

	return multiError(merrs)
}

func Combine(errs ...error) error { return Wrap(errs) }
