package reporter

import "io"

type ReportOptions struct {
	Tags map[string]interface{}

	Depth int
}

type Reporter interface {
	io.Closer

	Report(error, ReportOptions)
}
