package errors

import "github.com/upfluence/errors/stats"

func WithStatus(err error, status string) error {
	return stats.WithStatus(err, status)
}
