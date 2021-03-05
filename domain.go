package errors

import "github.com/upfluence/errors/domain"

func WithDomain(err error, d string) error {
	return domain.WithDomain(err, domain.Domain(d))
}
