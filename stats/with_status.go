package stats

type withStatus struct {
	cause  error
	status string
}

func (ws *withStatus) Error() string  { return ws.cause.Error() }
func (ws *withStatus) Unwrap() error  { return ws.cause }
func (ws *withStatus) Cause() error   { return ws.cause }
func (ws *withStatus) Status() string { return ws.status }

func (ws *withStatus) Tags() map[string]interface{} {
	return map[string]interface{}{"status": ws.status}
}

func WithStatus(err error, status string) error {
	if err == nil {
		return nil
	}

	return &withStatus{cause: err, status: status}
}
