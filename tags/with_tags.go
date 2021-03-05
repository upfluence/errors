package tags

type withTags struct {
	cause error
	tags  map[string]interface{}
}

func (ws *withTags) Error() string { return ws.cause.Error() }
func (ws *withTags) Unwrap() error { return ws.cause }
func (ws *withTags) Cause() error  { return ws.cause }

func (ws *withTags) Tags() map[string]interface{} {
	return ws.tags
}

func WithTags(err error, tags map[string]interface{}) error {
	if err == nil {
		return nil
	}

	return &withTags{cause: err, tags: tags}
}
