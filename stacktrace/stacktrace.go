package stacktrace

import (
	"runtime"
	"strings"

	"github.com/upfluence/errors/base"
)

type Frame uintptr

func Caller(depth int) Frame {
	var callers [1]uintptr
	runtime.Callers(2+depth, callers[:])

	return Frame(callers[0])
}

func Stacktrace(depth, count int) []Frame {
	var (
		callers = make([]uintptr, count)
		n       = runtime.Callers(2+depth, callers)
	)

	fs := make([]Frame, n)

	for i := 0; i < n; i++ {
		fs[i] = Frame(callers[i])
	}

	return fs

}

func (f Frame) Location() (string, string, int) {
	fr, _ := runtime.CallersFrames([]uintptr{uintptr(f)}).Next()

	return fr.Function, fr.File, fr.Line
}

func GetFrames(err error) []Frame {
	var fs []Frame

	for {
		if err == nil {
			break
		}

		switch ferr := err.(type) {
		case interface{ Frame() Frame }:
			fs = append(fs, ferr.Frame())
		case interface{ Frames() []Frame }:
			fs = append(fs, ferr.Frames()...)
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
