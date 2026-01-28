package stats

import (
	"errors"
	"fmt"
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
		"https://api.morisawafonts.com/webfont/v1/stats/pv/domains",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "2025-08", req.URL.Query().Get("from"))
			assert.Equal(t, "2025-09", req.URL.Query().Get("to"))

			switch req.URL.Query().Get(pager.Cursor) {
			case "":
				return httpmock.NewJsonResponse(http.StatusOK, pager.Page[*DomainsListResult, *DomainsListMetadata]{
					Result: []*DomainsListResult{
						{Domain: "1.example.com", Value: 100},
						{Domain: "2.example.com", Value: 200},
					},
					Meta: &DomainsListMetadata{
						Metadata: pager.Metadata{
							HasNext:    true,
							NextCursor: lo.ToPtr("cursor1"),
						},
						ProjectID: "project",
						From:      "2025-08",
						To:        "2025-09",
					},
				})
			case "cursor1":
				return httpmock.NewJsonResponse(http.StatusOK, pager.Page[*DomainsListResult, *DomainsListMetadata]{
					Result: []*DomainsListResult{
						{Domain: "3.example.com", Value: 300},
						{Domain: "4.example.com", Value: 400},
					},
					Meta: &DomainsListMetadata{
						Metadata: pager.Metadata{
							HasNext: false,
						},
						ProjectID: "project",
						From:      "2025-08",
						To:        "2025-09",
					},
				})
			}
			return nil, errors.New("unreachable")
		},
	)

	c := clienttest.NewClient(t)
	domains := NewDomains(c)

	list := domains.List(&DomainsListInput{
		From: lo.ToPtr("2025-08"),
		To:   lo.ToPtr("2025-09"),
	})

	i := 0
	for item, err := range list.Iter(t.Context()) {
		i++

		assert.Equal(t, &pager.Item[*DomainsListResult, *DomainsListMetadata]{
			Value: &DomainsListResult{
				Domain: fmt.Sprintf("%d.example.com", i),
				Value:  i * 100,
			},
			Meta: &DomainsListMetadata{
				Metadata: pager.Metadata{
					HasNext:    i <= 2,
					NextCursor: lo.Ternary(i <= 2, lo.ToPtr(fmt.Sprintf("cursor%d", (i+1)/2)), nil),
				},
				ProjectID: "project",
				From:      "2025-08",
				To:        "2025-09",
			},
		}, item)
		assert.NoError(t, err)
	}
	assert.Equal(t, 4, i)
}
