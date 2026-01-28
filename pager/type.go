package pager

import (
	"fmt"
	"net/url"
)

const (
	// Limit is the query parameter name for page size limit.
	Limit = "limit"
	// Cursor is the query parameter name for pagination cursor.
	Cursor = "cursor"
)

// Input represents pagination parameters for API requests.
type Input struct {
	Limit  *int
	Cursor *string
}

// Values converts the Input struct to url.Values for use in HTTP requests.
func (i *Input) Values() url.Values {
	if i == nil {
		return nil
	}

	values := url.Values{}
	if i.Limit != nil {
		values.Set(Limit, fmt.Sprintf("%d", *i.Limit))
	}
	if i.Cursor != nil {
		values.Set(Cursor, *i.Cursor)
	}
	return values
}

// Page represents a single page of paginated results.
// T is the type of items in the result, M is the metadata type.
type Page[T any, M metadata] struct {
	Result []T
	Meta   M
}

// Item wraps a single result item with its associated metadata.
// T is the type of the value, M is the metadata type.
type Item[T any, M metadata] struct {
	Value T
	Meta  M
}

type metadata interface {
	HasNextPage() bool
	NextPageCursor() *string
}

var _ metadata = (*Metadata)(nil)

// Metadata provides standard pagination metadata implementation.
type Metadata struct {
	HasNext    bool    `json:"has_next"`
	NextCursor *string `json:"next_cursor"`
}

// HasNextPage returns true if there are more pages available.
func (m *Metadata) HasNextPage() bool {
	return m.HasNext
}

// NextPageCursor returns the cursor for the next page, or nil if none.
func (m *Metadata) NextPageCursor() *string {
	return m.NextCursor
}
