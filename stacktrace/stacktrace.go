package stacktrace

import (
	"runtime"

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
