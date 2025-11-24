# errors

A comprehensive Go error handling library that extends the standard library with rich context, stack traces, error chaining, and integration with error reporting services.

This package provides a drop-in replacement for `errors` (stdlib) and `github.com/pkg/errors` with additional features like stacktrace embedding, tagging, domains, multi-error handling, and more.

## Features

- **Error Creation**: Create new errors with automatic domain detection and stack traces
- **Error Wrapping**: Add context to errors while preserving the original error chain
- **Stack Traces**: Automatic stack frame capture for debugging
- **Domains**: Categorize errors by domain (package-based or custom)
- **Tags**: Attach structured key-value metadata to errors
- **Status**: Associate status information with errors
- **Multi-Error Support**: Combine multiple errors into a single error
- **Opaque Errors**: Hide internal error details from external callers
- **Secondary Errors**: Attach additional contextual errors
- **Error Reporting**: Built-in support for Sentry and custom reporters
- **Standard Library Compatible**: Works with `errors.Is`, `errors.As`, and `errors.Unwrap`

## Installation

```bash
go get github.com/upfluence/errors
```

## Quick Start

```go
import "github.com/upfluence/errors"

// Create a new error
err := errors.New("configuration file not found")

// Create a formatted error
err := errors.Newf("invalid user ID: %d", userID)

// Wrap an existing error
file, err := os.Open("config.yaml")
if err != nil {
    return errors.Wrap(err, "failed to open configuration file")
}

// Wrap with formatting
data, err := readFile(path)
if err != nil {
    return errors.Wrapf(err, "failed to read file %s", path)
}
```

## Core Functions

### Creating Errors

**`New(msg string) error`**

Creates a new error with automatic domain detection, stack trace, and opaque wrapping.

```go
func loadConfig() error {
    if !fileExists("config.yaml") {
        return errors.New("configuration file not found")
    }
    return nil
}
```

**`Newf(msg string, args ...interface{}) error`**

Creates a new formatted error with automatic domain detection, stack trace, and opaque wrapping.

```go
func getUser(id int) error {
    if id < 0 {
        return errors.Newf("invalid user ID: %d", id)
    }
    return nil
}
```

### Wrapping Errors

**`Wrap(err error, msg string) error`**

Wraps an error with additional context and a stack frame.

```go
file, err := os.Open("config.yaml")
if err != nil {
    return errors.Wrap(err, "failed to open configuration file")
}
```

**`Wrapf(err error, msg string, args ...interface{}) error`**

Wraps an error with formatted context and a stack frame.

```go
data, err := readFile(path)
if err != nil {
    return errors.Wrapf(err, "failed to read file %s", path)
}
```

### Error Inspection

**`Cause(err error) error`**

Returns the root cause of an error by recursively unwrapping it.

```go
err := errors.New("root cause")
wrapped := errors.Wrap(err, "additional context")
root := errors.Cause(wrapped) // returns "root cause"
```

**`Unwrap(err error) error`**

Returns the next error in the chain (one level).

```go
err := errors.New("root")
wrapped := errors.Wrap(err, "wrapper")
unwrapped := errors.Unwrap(wrapped) // returns err
```

**`Is(err, target error) bool`**

Reports whether any error in the chain matches the target.

```go
if errors.Is(err, io.EOF) {
    fmt.Println("End of file reached")
}
```

**`As(err error, target interface{}) bool`**

Finds the first error in the chain matching the target type.

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Failed at path:", pathErr.Path)
}
```

**`IsTimeout(err error) bool`**

Checks if any error in the chain implements `Timeout() bool` and returns true.

```go
if errors.IsTimeout(err) {
    fmt.Println("Operation timed out, retrying...")
}
```

**`IsOfType[T error](err error) bool`** (Go 1.18+)

Reports whether any error in the chain matches the generic type T. This is a type-safe alternative to using `As` when you only need to check for the presence of a type without extracting it.

```go
type MyError struct{}
func (e MyError) Error() string { return "my error" }

err := errors.Wrap(MyError{}, "wrapped error")
if errors.IsOfType[MyError](err) {
    // MyError was found in the error chain
}
```

**`AsType[T error](err error) (T, bool)`** (Go 1.18+)

Attempts to convert an error to the generic type T by traversing the error chain. Returns the converted error and true if successful, or a zero value and false otherwise. This is a type-safe wrapper around `errors.As`.

```go
type MyError struct {
    Code int
}
func (e MyError) Error() string { return fmt.Sprintf("error %d", e.Code) }

err := errors.Wrap(MyError{Code: 404}, "not found")
if myErr, ok := errors.AsType[MyError](err); ok {
    fmt.Println(myErr.Code) // prints: 404
}
```

## Enriching Errors

### Stack Traces

**`WithStack(err error) error`**

Adds a stack frame at the current location.

```go
err := someExternalLibrary()
if err != nil {
    return errors.WithStack(err)
}
```

**`WithStack2[T any](v T, err error) (T, error)`** (Go 1.18+)

Adds a stack frame at the current location, while passing through a return value. This is useful for adding stack traces to errors from functions that return multiple values (value, error) in a single line.

If err is nil, returns (v, nil) unchanged. If err is not nil, returns (v, err_with_stack).

```go
// Instead of this:
result, err := externalLib.DoSomething()
if err != nil {
    return result, errors.WithStack(err)
}
return result, nil

// You can write this:
return errors.WithStack2(externalLib.DoSomething())
```

**`WithFrame(err error, depth int) error`**

Adds a stack frame at a specific depth in the call stack.

```go
func wrapError(err error) error {
    return errors.WithFrame(err, 1) // Skip current frame
}
```

### Domains

**`WithDomain(err error, domain string) error`**

Attaches a domain for error categorization.

```go
err := doSomething()
if err != nil {
    return errors.WithDomain(err, "database")
}
```

### Tags

**`WithTags(err error, tags map[string]interface{}) error`**

Attaches structured metadata to an error.

```go
err := fetchUser(userID)
if err != nil {
    return errors.WithTags(err, map[string]interface{}{
        "user_id": userID,
        "operation": "fetch",
    })
}
```

### Status

**`WithStatus(err error, status string) error`**

Attaches a status string to an error.

```go
err := processRequest()
if err != nil {
    return errors.WithStatus(err, "failed")
}
```

### Secondary Errors

**`WithSecondaryError(err error, additionalErr error) error`**

Attaches an additional error for context.

```go
err := saveToDatabase(data)
if err != nil {
    cacheErr := invalidateCache()
    return errors.WithSecondaryError(err, cacheErr)
}
```

## Multi-Error Support

**`Join(errs ...error) error`**

Combines multiple errors (compatible with Go 1.20+ `errors.Join`).

```go
err := errors.Join(
    errors.New("first error"),
    errors.New("second error"),
)
```

**`Combine(errs ...error) error`**

Alternative to `Join` with the same behavior.

```go
err := errors.Combine(
    validateName(),
    validateEmail(),
    validateAge(),
)
```

**`WrapErrors(errs []error) error`**

Combines a slice of errors.

```go
var errs []error
if err := validateName(); err != nil {
    errs = append(errs, err)
}
if err := validateEmail(); err != nil {
    errs = append(errs, err)
}
if len(errs) > 0 {
    return errors.WrapErrors(errs)
}
```

## Opaque Errors

**`Opaque(err error) error`**

Prevents unwrapping to hide internal implementation details.

```go
err := doInternalOperation()
if err != nil {
    return errors.Opaque(err) // Hide internal errors from callers
}
```

## Error Reporting

The package includes integration with error reporting services like Sentry.

### Reporter Interface

```go
type Reporter interface {
    io.Closer
    Report(error, ReportOptions)
}
```

### Sentry Integration

```go
import "github.com/upfluence/errors/reporter/sentry"

reporter, err := sentry.NewReporter(sentry.Options{
    DSN: "your-sentry-dsn",
})
if err != nil {
    log.Fatal(err)
}
defer reporter.Close()

reporter.Report(err, reporter.ReportOptions{
    Tags: map[string]interface{}{
        "environment": "production",
        "user_id": userID,
    },
})
```

### Standard Tag Keys

The package provides standard tag keys for common error metadata:

```go
reporter.TransactionKey       // "transaction"
reporter.DomainKey            // "domain"
reporter.UserEmailKey         // "user.email"
reporter.UserIDKey            // "user.id"
reporter.RemoteIP             // "remote.ip"
reporter.HTTPRequestPathKey   // "http.request.path"
reporter.HTTPRequestMethodKey // "http.request.method"
// ... and more
```

## Testing

The package includes testing utilities in the `errtest` subpackage for writing error-related tests.
