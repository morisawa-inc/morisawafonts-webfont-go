package option

import (
	"net/http"
	"net/url"
	"time"
)

type ClientOptions struct {
	APIToken   string
	BaseURL    *url.URL
	HTTPClient *http.Client
	Timeout    time.Duration
	Retry      int
}

func NewClientOptions(options ...Option) *ClientOptions {
	baseURL, _ := url.Parse(DefaultBaseURL)

	o := &ClientOptions{
		BaseURL: baseURL,
		Timeout: DefaultTimeout,
		Retry:   DefaultRetry,
	}
	for _, option := range options {
		option(o)
	}
	return o
}

func (o ClientOptions) Merge(options ...Option) *ClientOptions {
	merged := &o
	for _, option := range options {
		option(merged)
	}
	return merged
}
