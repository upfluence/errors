package errors

import "github.com/upfluence/errors/domain"

// WithDomain attaches a domain string to the error for categorization purposes.
func WithDomain(err error, d string) error {
	return domain.WithDomain(err, domain.Domain(d))
}
