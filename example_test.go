package errors_test

import (
	"fmt"
	"os"

	"github.com/upfluence/errors"
)

func ExampleNew() {
	err := errors.New("configuration file not found")
	fmt.Println(err)
	// Output: configuration file not found
}

func ExampleNewf() {
	id := -5
	err := errors.Newf("invalid user ID: %d", id)
	fmt.Println(err)
	// Output: invalid user ID: -5
}

func ExampleWrap() {
	// Simulate an error from os.Open
	baseErr := &os.PathError{Op: "open", Path: "config.yaml", Err: os.ErrNotExist}
	err := errors.Wrap(baseErr, "failed to open configuration file")
	fmt.Println(err)
	// Output: failed to open configuration file: open config.yaml: file does not exist
}

func ExampleWrapf() {
	// Simulate an error from reading a file
	path := "/tmp/data.txt"
	baseErr := &os.PathError{Op: "read", Path: path, Err: os.ErrPermission}
	err := errors.Wrapf(baseErr, "failed to read file %s", path)
	fmt.Println(err)
	// Output: failed to read file /tmp/data.txt: read /tmp/data.txt: permission denied
}
