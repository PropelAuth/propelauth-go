package models

import (
	"github.com/google/uuid"
)

// return types

type OrgMetadata struct {
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
}

type OrgList struct {
	TotalOrgs      int           `json:"total_orgs"`
	CurrentPage    int           `json:"current_page"`
	PageSize       int           `json:"page_size"`
	HasMoreResults bool          `json:"has_more_results"`
	Orgs           []OrgMetadata `json:"orgs"`
}

// post types

type CreateOrg struct {
	Name string `json:"name"`
}

type UpdateOrg struct {
	Name string `json:"name"`
}

type OrgQueryParams struct {
	PageSize   *int    `json:"page_size,omitempty"`
	PageNumber *int    `json:"page_number,omitempty"`
	OrderBy    *string `json:"order_by,omitempty"`
}
