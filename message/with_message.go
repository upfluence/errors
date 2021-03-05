package message

import "fmt"

type withMessage struct {
	cause error

	fmt  string
	args []interface{}
}

func (wm *withMessage) Error() string {
	var msg = wm.fmt

	if len(wm.args) > 0 {
		msg = fmt.Sprintf(msg, wm.args...)
	}

	return msg + ": " + wm.cause.Error()
}

func (wm *withMessage) Unwrap() error       { return wm.cause }
func (wm *withMessage) Cause() error        { return wm.cause }
func (wm *withMessage) Args() []interface{} { return wm.args }

func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &withMessage{cause: err, fmt: msg}
}

func WithMessagef(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withMessage{cause: err, fmt: msg, args: args}
}
