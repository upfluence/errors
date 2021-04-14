package stacktrace

import (
	"runtime"
	"strings"

	"github.com/upfluence/errors/base"
)

type Frame struct {
	frames [3]uintptr
}

func Caller(depth int) Frame {
	var s Frame

	runtime.Callers(depth+1, s.frames[:])

	return s
}

func (f Frame) Location() (string, string, int) {
	frames := runtime.CallersFrames(f.frames[:])

	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}

	fr, ok := frames.Next()

	if !ok {
		return "", "", 0
	}

	return fr.Function, fr.File, fr.Line
}

func GetFrames(err error) []Frame {
	var fs []Frame

	for {
		if err == nil {
			break
		}

		if ferr, ok := err.(interface{ Frame() Frame }); ok {
			fs = append(fs, ferr.Frame())
		}

		err = base.UnwrapOnce(err)
	}

	return fs
}

func PackageName(name string) string {
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
