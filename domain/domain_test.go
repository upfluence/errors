package domain_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/domain"
)

func TestDomain(t *testing.T) {
	err := errors.New("foo")

	assert.Equal(
		t,
		domain.Domain("github.com/upfluence/errors/domain"),
		domain.GetDomain(err),
	)

	assert.Equal(
		t,
		domain.Domain("github.com/upfluence/errors/domain"),
		domain.GetDomain(fmt.Errorf("wrapped: %w", err)),
	)

	assert.Equal(
		t,
		domain.Domain("bar"),
		domain.GetDomain(domain.WithDomain(err, "bar")),
	)

	assert.Equal(
		t,
		domain.NoDomain,
		domain.GetDomain(fmt.Errorf("error")),
	)
}
