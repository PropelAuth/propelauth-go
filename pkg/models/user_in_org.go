package models

import (
	"github.com/google/uuid"
)

// post types

type UserInOrgQueryParams struct {
	PageSize    *int  `json:"page_size,omitempty"`
	PageNumber  *int  `json:"page_number,omitempty"`
	IncludeOrgs *bool `json:"include_orgs,omitempty"`
}

type AddUserToOrg struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
	Role   string    `json:"role"`
}
