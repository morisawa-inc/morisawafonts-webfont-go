// Package morisawafonts provides a client for the Morisawa Fonts web font API.
package morisawafonts

import (
	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
	"github.com/morisawa-inc/morisawafonts-webfont-go/resource/domain"
)

// Client is the Morisawa Fonts web font API client.
type Client struct {
	*client.Client

	Domains *domain.Domains
}

// New creates a new Morisawa Fonts web font API client with the given options.
func New(options ...option.Option) *Client {
	c := client.NewClient(options...)
	return &Client{
		Client:  c,
		Domains: domain.NewDomains(c),
	}
}
