// Package option provides configuration options for the client.
package option

import (
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.morisawafonts.com/webfont/v1"
	DefaultTimeout = 30 * time.Second
	DefaultRetry   = 2
)

type Option func(*ClientOptions)

// WithAPIToken sets the API token for authentication.
func WithAPIToken(token string) Option {
	return func(o *ClientOptions) {
		o.APIToken = token
	}
}

// WithBaseURL sets a custom base URL for API requests.
func WithBaseURL(base *url.URL) Option {
	return func(o *ClientOptions) {
		o.BaseURL = base
	}
}

// WithHTTPClient sets a custom HTTP client for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(o *ClientOptions) {
		o.HTTPClient = client
	}
}

// WithTimeout sets the request timeout duration.
//
// Default: 30 seconds
func WithTimeout(timeout time.Duration) Option {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// WithRetry sets the number of retry attempts for failed requests.
//
// Default: 2
func WithRetry(retry int) Option {
	return func(o *ClientOptions) {
		o.Retry = retry
	}
}
