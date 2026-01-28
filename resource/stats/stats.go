// Package stats provides statistics functionality.
package stats

import "github.com/morisawa-inc/morisawafonts-webfont-go/client"

// Stats provides access to statistics.
type Stats struct {
	client *client.Client
	PV     *PV
}

// NewStats creates a new Stats instance.
func NewStats(c *client.Client) *Stats {
	return &Stats{
		client: c,
		PV:     NewPV(c),
	}
}
