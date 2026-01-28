package stats

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/morisawa-inc/morisawafonts-webfont-go/internal/clienttest"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPV_Get(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/stats/pv",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "2025-08", req.URL.Query().Get("from"))
			assert.Equal(t, "2025-09", req.URL.Query().Get("to"))

			return httpmock.NewJsonResponse(http.StatusOK, &PVGetResponse{
				PV: &PVGetResult{
					Total: 1234,
				},
				Meta: &PVGetMetadata{
					ProjectID: "project",
					From:      "2025-08",
					To:        "2025-09",
				},
			})
		},
	)

	c := clienttest.NewClient(t)
	pv := NewPV(c)

	result, err := pv.Get(t.Context(), &PVGetInput{
		From: lo.ToPtr("2025-08"),
		To:   lo.ToPtr("2025-09"),
	})

	assert.Equal(t, &PVGetResponse{
		PV: &PVGetResult{
			Total: 1234,
		},
		Meta: &PVGetMetadata{
			ProjectID: "project",
			From:      "2025-08",
			To:        "2025-09",
		},
	}, result)
	assert.NoError(t, err)
}
