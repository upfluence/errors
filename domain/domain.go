package domain

import (
	"go/build"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/upfluence/errors/base"
)

const NoDomain = "unknown"

type Domain string

func PackageDomain() Domain {
	return PackageDomainAtDepth(1)
}

func PackageDomainAtDepth(depth int) Domain {
	_, f, _, _ := runtime.Caller(1 + depth)

	dir := filepath.Dir(f)

	for _, d := range build.Default.SrcDirs() {
		dir = strings.TrimPrefix(dir, d+"/")
	}

	return Domain(dir)
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
