// Package pager provides pagination utilities.
package pager

import (
	"context"
	"io"
	"iter"
	"net/url"

	"github.com/morisawa-inc/morisawafonts-webfont-go/client"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
)

// Pager provides pagination functionality for API responses.
// T represents the type of items being paginated, M represents the metadata type.
type Pager[T any, M metadata] struct {
	client   *client.Client
	path     string
	values   url.Values
	options  []option.Option
	nextPage bool
}

// NewPager creates a new pager instance for paginating through API results.
func NewPager[T any, M metadata](
	c *client.Client,
	path string,
	values url.Values,
	options ...option.Option,
) *Pager[T, M] {
	return &Pager[T, M]{
		client:   c,
		path:     path,
		values:   values,
		options:  options,
		nextPage: true,
	}
}

// HasNextPage returns true if there are more pages to fetch.
func (p *Pager[T, M]) HasNextPage() bool {
	return p.nextPage
}

// GetNextPage fetches the next page of results.
// Returns io.EOF when no more pages are available.
func (p *Pager[T, M]) GetNextPage(ctx context.Context) (*Page[T, M], error) {
	if !p.nextPage {
		return nil, io.EOF
	}

	var page Page[T, M]
	err := p.client.Get(ctx, p.path, p.values, &page, p.options...)
	if err != nil {
		return nil, err
	}

	nextPage := page.Meta.HasNextPage()
	nextCursor := page.Meta.NextPageCursor()
	if nextPage && nextCursor != nil {
		if p.values == nil {
			p.values = url.Values{}
		}
		p.values.Set(Cursor, *nextCursor)
	} else {
		p.nextPage = false
	}

	return &page, nil
}

// Iter returns an iterator that yields individual items from all pages.
// The iterator automatically handles pagination and stops when all pages are consumed.
func (p *Pager[T, M]) Iter(ctx context.Context) iter.Seq2[*Item[T, M], error] {
	return func(yield func(*Item[T, M], error) bool) {
		for p.nextPage {
			page, err := p.GetNextPage(ctx)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, value := range page.Result {
				item := &Item[T, M]{
					Value: value,
					Meta:  page.Meta,
				}
				if !yield(item, nil) {
					return
				}
			}
		}
	}
}
