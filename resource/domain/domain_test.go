package domain

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/morisawa-inc/morisawafonts-webfont-go/internal/clienttest"
	"github.com/morisawa-inc/morisawafonts-webfont-go/pager"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDomains_List(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/domains",
		func(req *http.Request) (*http.Response, error) {
			switch req.URL.Query().Get(pager.Cursor) {
			case "":
				return httpmock.NewJsonResponse(http.StatusOK, pager.Page[string, *ListMetadata]{
					Result: []string{"1.example.com", "2.example.com"},
					Meta: &ListMetadata{
						Metadata: pager.Metadata{
							HasNext:    true,
							NextCursor: lo.ToPtr("cursor1"),
						},
						ProjectID: "project",
					},
				})
			case "cursor1":
				return httpmock.NewJsonResponse(http.StatusOK, pager.Page[string, *ListMetadata]{
					Result: []string{"3.example.com", "4.example.com"},
					Meta: &ListMetadata{
						Metadata: pager.Metadata{
							HasNext:    true,
							NextCursor: lo.ToPtr("cursor2"),
						},
						ProjectID: "project",
					},
				})
			case "cursor2":
				return httpmock.NewJsonResponse(http.StatusOK, pager.Page[string, *ListMetadata]{
					Result: []string{"5.example.com", "6.example.com"},
					Meta: &ListMetadata{
						Metadata: pager.Metadata{
							HasNext: false,
						},
						ProjectID: "project",
					},
				})
			}
			return nil, errors.New("unreachable")
		},
	)

	c := clienttest.NewClient(t)
	domains := NewDomains(c)

	list := domains.List(nil)

	i := 0
	for item, err := range list.Iter(t.Context()) {
		i++

		assert.Equal(t, &pager.Item[string, *ListMetadata]{
			Value: fmt.Sprintf("%d.example.com", i),
			Meta: &ListMetadata{
				Metadata: pager.Metadata{
					HasNext:    i <= 4,
					NextCursor: lo.Ternary(i <= 4, lo.ToPtr(fmt.Sprintf("cursor%d", (i+1)/2)), nil),
				},
				ProjectID: "project",
			},
		}, item)
		assert.NoError(t, err)
	}
	assert.Equal(t, 6, i)
}

func TestDomains_Add(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.morisawafonts.com/webfont/v1/domains",
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.JSONEq(t, `{"domains": ["example.com", "example.net"]}`, string(body))

			return httpmock.NewJsonResponse(http.StatusOK, &AddResult{
				Domains: []string{"example.com", "example.net"},
			})
		},
	)

	c := clienttest.NewClient(t)
	domains := NewDomains(c)

	result, err := domains.Add(t.Context(), []string{"example.com", "example.net"})

	assert.Equal(t, &AddResult{
		Domains: []string{"example.com", "example.net"},
	}, result)
	assert.NoError(t, err)
}

func TestDomains_Delete(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodDelete,
		"https://api.morisawafonts.com/webfont/v1/domains",
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.JSONEq(t, `{"domains": ["example.com", "example.net"]}`, string(body))

			return httpmock.NewBytesResponse(http.StatusNoContent, nil), nil
		},
	)

	c := clienttest.NewClient(t)
	domains := NewDomains(c)

	err := domains.Delete(t.Context(), []string{"example.com", "example.net"})

	assert.NoError(t, err)
}
