// Package stacktrace provides utilities for capturing and managing stack traces.
//
// This package allows capturing stack frame information at error creation time,
// which can be attached to errors and later extracted for debugging or reporting.
// It supports both single frame and full stacktrace capture, and provides utilities
// for extracting package names from function names.
package stacktrace

import (
	"runtime"
	"strings"

	"github.com/upfluence/errors/base"
)

// Frame represents a single program counter location in a stack trace.
type Frame uintptr

// Caller captures a single stack frame at the specified depth.
// depth=0 returns the caller's frame, depth=1 returns the caller's caller's frame, etc.
func Caller(depth int) Frame {
	var callers [1]uintptr

	runtime.Callers(2+depth, callers[:])

	return Frame(callers[0])
}

// Stacktrace captures multiple stack frames starting at the specified depth.
// Returns up to count frames from the call stack.
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

// Location returns the function name, file path, and line number for this frame.
func (f Frame) Location() (string, string, int) {
	fr, _ := runtime.CallersFrames([]uintptr{uintptr(f)}).Next()

	return fr.Function, fr.File, fr.Line
}

// GetFrames extracts all stack frames from an error by traversing the error chain.
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

// PackageName extracts the package path from a fully qualified function name.
// Returns an empty string for compiler-generated symbols.
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
