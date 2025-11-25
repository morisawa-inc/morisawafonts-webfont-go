// Package domain provides domain management functionality.
package domain

import (
	"context"

	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
	"github.com/morisawa-inc/morisawafonts-webfont-go/pager"
)

// Domains provides domain management operations.
type Domains struct {
	client *client.Client
}

// NewDomains creates a new Domains instance.
func NewDomains(c *client.Client) *Domains {
	return &Domains{
		client: c,
	}
}

// List returns a paginated list of domains.
func (d *Domains) List(
	input *ListInput,
	options ...option.Option,
) *pager.Pager[string, *ListMetadata] {
	return pager.NewPager[string, *ListMetadata](d.client, "/domains", input.Values(), options...)
}

// Add adds the specified domains.
func (d *Domains) Add(
	ctx context.Context,
	domains []string,
	options ...option.Option,
) (*AddResult, error) {
	body := map[string]any{
		"domains": domains,
	}

	var result AddResult
	err := d.client.Post(ctx, "/domains", body, &result, options...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes the specified domains.
func (d *Domains) Delete(
	ctx context.Context,
	domains []string,
	options ...option.Option,
) error {
	body := map[string]any{
		"domains": domains,
	}

	return d.client.Delete(ctx, "/domains", body, options...)
}
