package stacktrace

type withFrame struct {
	cause error
	frame Frame
}

func (wf *withFrame) Error() string { return wf.cause.Error() }
func (wf *withFrame) Unwrap() error { return wf.cause }
func (wf *withFrame) Cause() error  { return wf.cause }
func (wf *withFrame) Frame() Frame  { return wf.frame }

func WithFrame(err error, depth int) error {
	if err == nil {
		return nil
	}

	return &withFrame{cause: err, frame: Caller(depth + 1)}
}
