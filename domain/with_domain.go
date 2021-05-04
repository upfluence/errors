package domain

import "errors"

type withDomain struct {
	cause error

	domain Domain
}

func (ws *withDomain) Error() string  { return ws.cause.Error() }
func (ws *withDomain) Unwrap() error  { return ws.cause }
func (ws *withDomain) Cause() error   { return ws.cause }
func (ws *withDomain) Domain() Domain { return ws.domain }

func (ws *withDomain) Tags() map[string]interface{} {
	return map[string]interface{}{"domain": string(ws.domain)}
}

func WithDomain(err error, d Domain) error {
	if err == nil {
		return nil
	}

	return &withDomain{cause: err, domain: d}
}

func New(msg string) error {
	return WithDomain(errors.New(msg), PackageDomain())
}
