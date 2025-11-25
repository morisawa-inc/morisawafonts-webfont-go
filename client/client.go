// Package client provides a low-level HTTP client.
package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
	"resty.dev/v3"
)

// Client provides a low-level HTTP client for API communication.
type Client struct {
	options *option.ClientOptions
	resty   *resty.Client
}

// NewClient creates a new HTTP client with the specified options.
func NewClient(
	options ...option.Option,
) *Client {
	o := option.NewClientOptions(options...)

	return &Client{
		options: o,
		resty:   setupResty(o),
	}
}

// Close closes the underlying HTTP client and releases resources.
func (c *Client) Close() error {
	return c.resty.Close()
}

// Get performs a GET request to the specified path with query parameters.
func (c *Client) Get(
	ctx context.Context,
	path string,
	values url.Values,
	result any,
	options ...option.Option,
) error {
	return c.request(ctx, http.MethodGet, path, values, nil, result, options...)
}

// Post performs a POST request to the specified path with a request body.
func (c *Client) Post(
	ctx context.Context,
	path string,
	body any,
	result any,
	options ...option.Option,
) error {
	return c.request(ctx, http.MethodPost, path, nil, body, result, options...)
}

// Delete performs a DELETE request to the specified path with a request body.
func (c *Client) Delete(
	ctx context.Context,
	path string,
	body any,
	options ...option.Option,
) error {
	return c.request(ctx, http.MethodDelete, path, nil, body, nil, options...)
}

func (c *Client) request(
	ctx context.Context,
	method string,
	path string,
	values url.Values,
	body any,
	result any,
	options ...option.Option,
) error {
	o := c.options.Merge(options...)
	if o.APIToken == "" {
		return ErrNoAPIToken
	}

	r := c.resty

	// temporary client
	if c.options.HTTPClient != o.HTTPClient {
		r = setupResty(o)
		defer func() {
			_ = r.Close()
		}()
	}

	response, err := r.R().
		SetContext(ctx).
		SetAuthToken(o.APIToken).
		SetQueryParamsFromValues(values).
		SetBody(body).
		SetResult(result).
		SetTimeout(o.Timeout).
		SetRetryCount(o.Retry).
		Execute(method, o.BaseURL.JoinPath(path).String())
	if err != nil {
		return err
	}
	if response.IsError() {
		return NewAPIError(response)
	}
	return nil
}

func setupResty(o *option.ClientOptions) *resty.Client {
	var r *resty.Client
	if o.HTTPClient != nil {
		r = resty.NewWithClient(o.HTTPClient)
	} else {
		r = resty.New()
	}

	r.SetHeader("user-agent", getUserAgent()).
		SetAllowMethodDeletePayload(true).
		SetLogger(&logger{})

	return r
}
