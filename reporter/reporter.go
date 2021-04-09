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
	HTTPRequestMethodKey = "http.request.method"

	HTTPRequestHeaderKeyPrefix      = "http.request.headers."
	HTTPRequestQueryValuesKeyPrefix = "http.request.query_values."

	ThriftRequestMethodKey  = "thrift.request.method"
	ThriftRequestServiceKey = "thrift.request.service"
	ThriftRequestCallerKey  = "thrift.request.caller"
)

type ReportOptions struct {
	Tags map[string]interface{}

	Depth int
}

type Reporter interface {
	io.Closer

	Report(error, ReportOptions)
}
