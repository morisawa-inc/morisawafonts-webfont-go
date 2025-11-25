package pager

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/morisawa-inc/morisawafonts-webfont-go/internal/clienttest"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func setupMock(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/pager",
		func(req *http.Request) (*http.Response, error) {
			switch req.URL.Query().Get(Cursor) {
			case "":
				return httpmock.NewJsonResponse(http.StatusOK, Page[int, *Metadata]{
					Result: []int{1, 2, 3},
					Meta: &Metadata{
						HasNext:    true,
						NextCursor: lo.ToPtr("cursor1"),
					},
				})
			case "cursor1":
				return httpmock.NewJsonResponse(http.StatusOK, Page[int, *Metadata]{
					Result: []int{4, 5, 6},
					Meta: &Metadata{
						HasNext:    true,
						NextCursor: lo.ToPtr("cursor2"),
					},
				})
			case "cursor2":
				return httpmock.NewJsonResponse(http.StatusOK, Page[int, *Metadata]{
					Result: []int{7, 8, 9},
					Meta: &Metadata{
						HasNext: false,
					},
				})
			}
			return nil, errors.New("unreachable")
		},
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/empty",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, Page[int, *Metadata]{
				Result: []int{},
				Meta: &Metadata{
					HasNext: false,
				},
			})
		},
	)
}

func TestPager_GetNextPage(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantResults [][]int
	}{
		{
			"pager",
			"/pager",
			[][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			"empty",
			"/empty",
			[][]int{
				{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupMock(t)

			c := clienttest.NewClient(t)
			pager := NewPager[int, *Metadata](c, tt.path, nil)

			assert.True(t, pager.HasNextPage())

			for i, want := range tt.wantResults {
				page, err := pager.GetNextPage(t.Context())
				assert.Equal(
					t,
					&Page[int, *Metadata]{
						Result: want,
						Meta: &Metadata{
							HasNext:    i < len(tt.wantResults)-1,
							NextCursor: lo.Ternary(i < len(tt.wantResults)-1, lo.ToPtr(fmt.Sprintf("cursor%d", i+1)), nil),
						},
					},
					page,
				)
				assert.NoError(t, err)
			}

			assert.False(t, pager.HasNextPage())

			// overflow
			page, err := pager.GetNextPage(t.Context())
			assert.Nil(t, page)
			assert.ErrorIs(t, err, io.EOF)
		})
	}
}

func TestPager_Iter(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantValues []int
		perPage    int
	}{
		{
			"pager",
			"/pager",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			3,
		},
		{
			"empty",
			"/empty",
			nil,
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupMock(t)

			c := clienttest.NewClient(t)
			pager := NewPager[int, *Metadata](c, tt.path, nil)

			assert.True(t, pager.HasNextPage())

			i := 0
			for item, err := range pager.Iter(t.Context()) {
				i++
				assert.Equal(
					t,
					&Item[int, *Metadata]{
						Value: tt.wantValues[i-1],
						Meta: &Metadata{
							HasNext:    i <= tt.perPage*(len(tt.wantValues)/tt.perPage-1),
							NextCursor: lo.Ternary(i <= tt.perPage*(len(tt.wantValues)/tt.perPage-1), lo.ToPtr(fmt.Sprintf("cursor%d", (i+tt.perPage-1)/tt.perPage)), nil),
						},
					},
					item,
				)
				assert.NoError(t, err)
			}
			assert.Len(t, tt.wantValues, i)

			assert.False(t, pager.HasNextPage())
		})
	}
}
