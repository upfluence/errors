package domain

import (
	"strings"

	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/stacktrace"
)

const NoDomain = Domain("unknown")

type Domain string

func PackageDomain() Domain {
	return PackageDomainAtDepth(1)
}

func PackageDomainAtDepth(depth int) Domain {
	var fn, _, _ = stacktrace.Caller(1 + depth).Location()

	return Domain(packageName(fn))
}

func GetDomain(err error) Domain {
	for {
		if wd, ok := err.(*withDomain); ok {
			return wd.domain
		}

		err = base.UnwrapOnce(err)

		if err == nil {
			break
		}
	}

	return NoDomain
}

func packageName(name string) string {
	// A prefix of "type." and "go." is a compiler-generated symbol that doesn't belong to any package.
	// See variable reservedimports in cmd/compile/internal/gc/subr.go
	if strings.HasPrefix(name, "go.") || strings.HasPrefix(name, "type.") {
		return ""
	}

	pathend := strings.LastIndex(name, "/")
	if pathend < 0 {
		pathend = 0
	}

	if i := strings.Index(name[pathend:], "."); i != -1 {
		return name[:pathend+i]
	}
	return ""
}
