package errors

import (
	"errors"

	"github.com/upfluence/errors/domain"
	"github.com/upfluence/errors/message"
	"github.com/upfluence/errors/opaque"
)

func New(msg string) error {
	return opaque.Opaque(
		domain.WithDomain(
			WithFrame(errors.New(msg), 1),
			domain.PackageDomainAtDepth(1),
		),
	)
}

func Wrap(err error, msg string) error {
	return WithFrame(message.WithMessage(err, msg), 1)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return WithFrame(message.WithMessagef(err, msg, args...), 1)
}
