package stats

import (
	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
	"github.com/morisawa-inc/morisawafonts-webfont-go/pager"
)

// Domains retrieves page view statistics by domain.
type Domains struct {
	client *client.Client
}

// NewDomains creates a new Domains instance.
func NewDomains(c *client.Client) *Domains {
	return &Domains{
		client: c,
	}
}

// List returns a paginated list of page view statistics by domain.
func (d *Domains) List(
	input *DomainsListInput,
	options ...option.Option,
) *pager.Pager[*DomainsListResult, *DomainsListMetadata] {
	return pager.NewPager[*DomainsListResult, *DomainsListMetadata](d.client, "/stats/pv/domains", input.Values(), options...)
}
