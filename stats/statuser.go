package stats

import (
	"fmt"

	"github.com/upfluence/errors/base"
)

type Statuser interface {
	Status(error) string
}

type ExtractStatusOption func(*extractStatusOptions)

type extractStatusOptions struct {
	successStatus    string
	fallbackStatuser Statuser
}

var defaultStatusOptions = extractStatusOptions{
	successStatus:    "success",
	fallbackStatuser: defaultStatuser{},
}

type defaultStatuser struct{}

func (defaultStatuser) Status(err error) string {
	switch t := fmt.Sprintf("%T", err); t {
	case "*errors.errorString", "*errors.fundamental":
		return err.Error()
	default:
		return t
	}
}

func GetStatus(err error, opts ...ExtractStatusOption) string {
	var o = defaultStatusOptions

	for _, opt := range opts {
		opt(&o)
	}

	if err == nil {
		return o.successStatus
	}

	for {
		ws, ok := err.(*withStatus)

		if ok {
			return ws.status
		}

		cause := base.UnwrapOnce(err)

		if cause == nil {
			break
		}

		err = cause
	}

	return o.fallbackStatuser.Status(err)
}
