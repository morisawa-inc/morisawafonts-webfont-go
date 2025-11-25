package client

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestNewAPIError(t *testing.T) {
	type args struct {
		response *resty.Response
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantStatus     string
		wantError      string
	}{
		{
			"raw=empty",
			args{
				&resty.Response{
					Request: &resty.Request{
						URL:    "http://example.com/",
						Method: http.MethodGet,
					},
					RawResponse: &http.Response{
						StatusCode: 404,
						Status:     "404 Not Found",
					},
					Body: io.NopCloser(strings.NewReader("")),
				},
			},
			404,
			"404 Not Found",
			"api error: 404 Not Found: GET http://example.com/",
		},
		{
			"raw=text",
			args{
				&resty.Response{
					Request: &resty.Request{
						URL:    "http://example.com/",
						Method: http.MethodGet,
					},
					RawResponse: &http.Response{
						StatusCode: 404,
						Status:     "404 Not Found",
					},
					Body: io.NopCloser(strings.NewReader("this is a test")),
				},
			},
			404,
			"404 Not Found",
			"api error: 404 Not Found: GET http://example.com/: this is a test",
		},
		{
			"raw=json",
			args{
				&resty.Response{
					Request: &resty.Request{
						URL:    "http://example.com/",
						Method: http.MethodGet,
					},
					RawResponse: &http.Response{
						StatusCode: 404,
						Status:     "404 Not Found",
					},
					Body: io.NopCloser(strings.NewReader(`{"message": "this is a test"}`)),
				},
			},
			404,
			"404 Not Found",
			"api error: 404 Not Found: GET http://example.com/: this is a test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewAPIError(tt.args.response)

			assert.Equal(t, tt.wantStatusCode, got.StatusCode)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.wantError, got.Error())
		})
	}
}
