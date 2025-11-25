package domain

import "github.com/morisawa-inc/morisawafonts-webfont-go/pager"

type ListInput = pager.Input

type ListMetadata struct {
	pager.Metadata

	ProjectID string `json:"project_id"`
}

type AddResult struct {
	Domains []string `json:"domains"`
}
