package errors

import (
	"errors"
	"fmt"

	"github.com/upfluence/errors/domain"
	"github.com/upfluence/errors/message"
	"github.com/upfluence/errors/opaque"
)

// New creates a new error with the given message. The error is wrapped with a
// stack frame, domain derived from the calling package, and made opaque.
func New(msg string) error {
	return opaque.Opaque(
		domain.WithDomain(
			WithFrame(errors.New(msg), 1),
			domain.PackageDomainAtDepth(1),
		),
	)
}

// Newf creates a new error with a formatted message. The error is wrapped with a
// stack frame, domain derived from the calling package, and made opaque.
func Newf(msg string, args ...interface{}) error {
	return opaque.Opaque(
		domain.WithDomain(
			WithFrame(fmt.Errorf(msg, args...), 1),
			domain.PackageDomainAtDepth(1),
		),
	)
}

// Wrap wraps an error with an additional message and stack frame.
func Wrap(err error, msg string) error {
	return WithFrame(message.WithMessage(err, msg), 1)
}

// Wrapf wraps an error with a formatted message and stack frame.
func Wrapf(err error, msg string, args ...interface{}) error {
	return WithFrame(message.WithMessagef(err, msg, args...), 1)
}
