package stats

import (
	"context"

	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
)

// PV retrieves page view statistics.
type PV struct {
	client  *client.Client
	Domains *Domains
}

// NewPV creates a new PV instance.
func NewPV(c *client.Client) *PV {
	return &PV{
		client:  c,
		Domains: NewDomains(c),
	}
}

// Get retrieves page view statistics for the project.
func (p *PV) Get(
	ctx context.Context,
	input *PVGetInput,
	options ...option.Option,
) (*PVGetResponse, error) {
	var result PVGetResponse
	err := p.client.Get(ctx, "/stats/pv", input.Values(), &result, options...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
