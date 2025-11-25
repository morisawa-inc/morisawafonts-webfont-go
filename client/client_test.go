package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
	"github.com/stretchr/testify/assert"
)

type testResult struct {
	Data string `json:"data"`
}

type testBody struct {
	Param int `json:"param"`
}

func TestClient_Get(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/get",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, url.Values{"param": {"test param"}}, req.URL.Query())

			return httpmock.NewJsonResponse(http.StatusOK, &testResult{Data: "some data"})
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	var result testResult
	err := c.Get(t.Context(), "/get", url.Values{"param": {"test param"}}, &result)

	assert.NoError(t, err)
	assert.Equal(t, "some data", result.Data)
}

func TestClient_Post(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.morisawafonts.com/webfont/v1/post",
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.JSONEq(t, `{"param": 42}`, string(body))

			return httpmock.NewJsonResponse(http.StatusOK, &testResult{Data: "some data"})
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	var result testResult
	err := c.Post(t.Context(), "/post", testBody{Param: 42}, &result)

	assert.NoError(t, err)
	assert.Equal(t, "some data", result.Data)
}

func TestClient_Delete(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodDelete,
		"https://api.morisawafonts.com/webfont/v1/delete",
		func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.JSONEq(t, `{"param": 42}`, string(body))

			return httpmock.NewBytesResponse(http.StatusNoContent, nil), nil
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	err := c.Delete(t.Context(), "/delete", testBody{Param: 42})

	assert.NoError(t, err)
}

func TestClient_error_noToken(t *testing.T) {
	httpmock.Activate(t)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
	)

	err := c.Get(t.Context(), "/get", nil, nil)

	assert.ErrorIs(t, err, ErrNoAPIToken)
}

func TestClient_error_api(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/get",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusBadRequest, map[string]string{"message": "error message"})
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	err := c.Get(t.Context(), "/get", nil, nil)

	var apiErr *APIError
	if assert.ErrorAs(t, err, &apiErr) {
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "400 Bad Request", apiErr.Status)
		assert.Equal(t, "api error: 400 Bad Request: GET https://api.morisawafonts.com/webfont/v1/get: error message", apiErr.Error())
	}
}

func TestClient_error_network(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/get",
		httpmock.ConnectionFailure,
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	err := c.Get(t.Context(), "/get", nil, nil)

	var urlErr *url.Error
	if assert.ErrorAs(t, err, &urlErr) {
		assert.ErrorIs(t, urlErr.Err, httpmock.NoResponderFound)
	}
}

func TestClient_error_timeout(t *testing.T) {
	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/get",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(500 * time.Millisecond)
			return nil, errors.New("unreachable")
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
		option.WithTimeout(100*time.Millisecond),
	)

	err := c.Get(t.Context(), "/get", nil, nil)

	var urlErr *url.Error
	if assert.ErrorAs(t, err, &urlErr) {
		assert.ErrorIs(t, urlErr.Err, context.DeadlineExceeded)
	}
}

func TestClient_retry(t *testing.T) {
	retry := 0

	httpmock.Activate(t)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.morisawafonts.com/webfont/v1/get",
		func(req *http.Request) (*http.Response, error) {
			retry++
			if retry < 3 {
				return httpmock.NewJsonResponse(http.StatusTooManyRequests, nil)
			}
			return httpmock.NewJsonResponse(http.StatusOK, &testResult{Data: "some data"})
		},
	)

	c := NewClient(
		option.WithHTTPClient(http.DefaultClient),
		option.WithAPIToken("test-token"),
	)

	var result testResult
	err := c.Get(t.Context(), "/get", nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "some data", result.Data)

	assert.Equal(t, 3, retry)
}
