package domain

import (
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

	return Domain(stacktrace.PackageName(fn))
}

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
