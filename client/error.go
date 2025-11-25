package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"resty.dev/v3"
)

var (
	// ErrNoAPIToken is returned when no API token is provided.
	ErrNoAPIToken = errors.New("api token is required")

	_ error = (*APIError)(nil)
)

// APIError represents an error response from the API.
type APIError struct {
	StatusCode int
	Status     string

	message string
}

// NewAPIError creates a new APIError from an HTTP response.
func NewAPIError(response *resty.Response) *APIError {
	message := formatErrorMessage(response)

	return &APIError{
		StatusCode: response.StatusCode(),
		Status:     response.Status(),
		message:    message,
	}
}

// Error returns the error message for the APIError.
func (err *APIError) Error() string {
	return err.message
}

func formatErrorMessage(response *resty.Response) string {
	message := fmt.Sprintf(
		"api error: %s: %s %s",
		response.Status(),
		response.Request.Method,
		response.Request.URL,
	)

	raw := response.String()
	if len(raw) > 0 {
		var result struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal([]byte(raw), &result)
		if err == nil && result.Message != "" {
			message = fmt.Sprintf("%s: %s", message, result.Message)
		} else {
			message = fmt.Sprintf("%s: %s", message, raw)
		}
	}

	return message
}
