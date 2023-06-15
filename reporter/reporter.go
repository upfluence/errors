package reporter

import "io"

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

var NopReporter Reporter = nopReporter{}

type ReportOptions struct {
	Tags map[string]interface{}

	Depth int
}

type Reporter interface {
	io.Closer

	Report(error, ReportOptions)
}

type nopReporter struct{}

func (nopReporter) Close() error                { return nil }
func (nopReporter) Report(error, ReportOptions) {}
