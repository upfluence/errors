// Package domain provides error domain classification and extraction.
//
// A domain identifies the package or module where an error originated,
// which is useful for categorizing and routing errors in large applications.
// The domain is automatically extracted from the call stack and can be
// attached to errors for later retrieval.
package domain

import (
	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/stacktrace"
)

// NoDomain is returned when an error has no associated domain.
const NoDomain = Domain("unknown")

// Domain represents the package or module where an error originated.
type Domain string

// PackageDomain returns the domain for the calling package.
func PackageDomain() Domain {
	return PackageDomainAtDepth(1)
}

// PackageDomainAtDepth returns the domain for the package at the specified
// call stack depth. depth=0 returns the caller's package, depth=1 returns
// the caller's caller's package, etc.
func PackageDomainAtDepth(depth int) Domain {
	var fn, _, _ = stacktrace.Caller(1 + depth).Location()

	return Domain(stacktrace.PackageName(fn))
}

// GetDomain extracts the domain from an error by traversing the error chain.
// Returns NoDomain if no domain is found.
func GetDomain(err error) Domain {
	for {
		if wd, ok := err.(interface{ Domain() Domain }); ok {
			return wd.Domain()
		}

		err = base.UnwrapOnce(err)

		if err == nil {
			break
		}
	}

	return NoDomain
}
