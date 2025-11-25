package client

import "resty.dev/v3"

var _ resty.Logger = (*logger)(nil)

type logger struct{}

func (l *logger) Errorf(_ string, _ ...any) {
	// TODO: not implemented
}

func (l *logger) Warnf(_ string, _ ...any) {
	// TODO: not implemented
}

func (l *logger) Debugf(_ string, _ ...any) {
	// TODO: not implemented
}
