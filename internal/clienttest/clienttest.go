// Package clienttest provides testing utilities for the client.
package clienttest

import (
	"net/http"
	"testing"

	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
)

func NewClient(t *testing.T, options ...option.Option) *client.Client {
	options = append(
		[]option.Option{
			option.WithHTTPClient(http.DefaultClient),
			option.WithAPIToken("test-token"),
		},
		options...,
	)

	c := client.NewClient(options...)
	t.Cleanup(func() {
		_ = c.Close()
	})
	return c
}
