package stacktrace

type withStacktrace struct {
	cause  error
	frames []Frame
}

func (ws *withStacktrace) Error() string   { return ws.cause.Error() }
func (ws *withStacktrace) Unwrap() error   { return ws.cause }
func (ws *withStacktrace) Cause() error    { return ws.cause }
func (ws *withStacktrace) Frames() []Frame { return ws.frames }

func WithStacktrace(err error, depth, count int) error {
	if err == nil {
		return nil
	}

	return &withStacktrace{cause: err, frames: Stacktrace(depth+1, count)}
}
