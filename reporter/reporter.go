// Package reporter provides an interface for error reporting to external services.
//
// This package defines the Reporter interface and common tag keys for structured
// error reporting. Reporters can send errors to services like Sentry with rich
// contextual information including request details, user information, and custom tags.
package reporter

import "io"

// Common tag keys for error reporting metadata.
const (
	TransactionKey = "transaction"
	DomainKey      = "domain"

	UserEmailKey = "user.email"
	UserIDKey    = "user.id"

	RemoteIP   = "remote.ip"
	RemotePort = "remote.port"

	HTTPRequestPathKey   = "http.request.path"
	HTTPRequestHostKey   = "http.request.host"
	HTTPRequestProtoKey  = "http.request.proto"
	HTTPRequestPortKey   = "http.request.port"
	HTTPRequestMethodKey = "http.request.method"
	HTTPRequestBodyKey   = "http.request.body"

	HTTPRequestHeaderKeyPrefix      = "http.request.headers."
	HTTPRequestQueryValuesKeyPrefix = "http.request.query_values."

	ThriftRequestMethodKey  = "thrift.request.method"
	ThriftRequestServiceKey = "thrift.request.service"
	ThriftRequestCallerKey  = "thrift.request.caller"
	ThriftRequestBodyKey    = "thrift.request.body"
)

// NopReporter is a Reporter implementation that does nothing.
var NopReporter Reporter = nopReporter{}

// ReportOptions contains options for reporting an error.
type ReportOptions struct {
	Tags map[string]interface{}

	Depth int
}

// Reporter is the interface for error reporting implementations.
type Reporter interface {
	io.Closer

	Report(error, ReportOptions)
}

type nopReporter struct{}

func (nopReporter) Close() error                { return nil }
func (nopReporter) Report(error, ReportOptions) {}
