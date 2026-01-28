package stats

import (
	"net/url"

	"github.com/morisawa-inc/morisawafonts-webfont-go/pager"
	"github.com/samber/lo"
)

type PVGetInput struct {
	From *string
	To   *string
}

func (i *PVGetInput) Values() url.Values {
	if i == nil {
		return nil
	}

	values := url.Values{}
	if i.From != nil {
		values.Set("from", *i.From)
	}
	if i.To != nil {
		values.Set("to", *i.To)
	}
	return values
}

type PVGetResponse struct {
	PV   *PVGetResult   `json:"pv"`
	Meta *PVGetMetadata `json:"meta"`
}

type PVGetResult struct {
	Total int `json:"total"`
}

type PVGetMetadata struct {
	ProjectID string `json:"project_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}

type DomainsListInput struct {
	pager.Input

	From   *string
	To     *string
	Domain *string
}

func (i *DomainsListInput) Values() url.Values {
	if i == nil {
		return nil
	}

	values := url.Values{}
	if i.From != nil {
		values.Set("from", *i.From)
	}
	if i.To != nil {
		values.Set("to", *i.To)
	}
	if i.Domain != nil {
		values.Set("domain", *i.Domain)
	}

	return lo.Assign(values, i.Input.Values())
}

type DomainsListResult struct {
	Domain string `json:"domain"`
	Value  int    `json:"value"`
}

type DomainsListMetadata struct {
	pager.Metadata

	ProjectID string `json:"project_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}
