package models

import (
	"github.com/google/uuid"
)

// return types

type UserInOrgQueryParams struct {
	PageSize    int  `json:"page_size"`
	PageNumber  int  `json:"page_number"`
	IncludeOrgs bool `json:"include_orgs"`
}

// post types

type AddUserToOrg struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
	Role   string    `json:"role"`
}
