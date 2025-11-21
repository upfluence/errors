package errors_test

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/upfluence/errors"
)

func ExampleCause() {
	rootErr := errors.New("root cause")
	wrapped := errors.Wrap(rootErr, "additional context")
	root := errors.Cause(wrapped)
	fmt.Println(root)
	// Output: root cause
}

func ExampleUnwrap() {
	rootErr := errors.New("root")
	wrapped := errors.Wrap(rootErr, "wrapper")
	unwrapped := errors.Unwrap(wrapped)
	fmt.Println(unwrapped)
	// Output: wrapper: root
}

func ExampleAs() {
	// Create a PathError
	pathErr := &os.PathError{Op: "open", Path: "/tmp/file.txt", Err: os.ErrNotExist}
	wrapped := errors.Wrap(pathErr, "failed to process file")

	// Use As to extract the PathError
	var targetErr *os.PathError
	if errors.As(wrapped, &targetErr) {
		fmt.Println("Failed at path:", targetErr.Path)
	}
	// Output: Failed at path: /tmp/file.txt
}

func ExampleIs() {
	err := errors.Wrap(io.EOF, "read operation failed")

	if errors.Is(err, io.EOF) {
		fmt.Println("End of file reached")
	}
	// Output: End of file reached
}

func ExampleIsTimeout() {
	// Create a timeout error
	timeoutErr := &net.DNSError{IsTimeout: true}
	wrapped := errors.Wrap(timeoutErr, "network operation failed")

	if errors.IsTimeout(wrapped) {
		fmt.Println("Operation timed out, retrying...")
	}
	// Output: Operation timed out, retrying...
}
