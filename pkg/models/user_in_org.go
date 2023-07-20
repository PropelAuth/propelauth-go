package models

import (
	"github.com/google/uuid"
)

// post types

// UserInOrgQueryParams is the information for querying a pageable list of users in an organization.
// If left blank, PageSize defaults to 10 and PageNumber defaults to 0. IncludeOrgs defaults to false,
// but if set to true will include all orgs each user is in.
type UserInOrgQueryParams struct {
	PageSize    *int  `json:"page_size,omitempty"`
	PageNumber  *int  `json:"page_number,omitempty"`
	IncludeOrgs *bool `json:"include_orgs,omitempty"`
}

// AddUserToOrg is the information needed to add a user to an organization. Role is just a string, but
// it has to match up to one of your defined roles, by default these are Owner, Admin, or Member, but
// they can be changed via your dashboard.
type AddUserToOrg struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
	Role   string    `json:"role"`
}

// InviteUserToOrg is the information needed to invite a new user to join an organization. Role is
// just a string, but it has to match up to one of your defined roles, by default these are Owner,
// Admin, or Member, but they can be changed via your dashboard.
type InviteUserToOrg struct {
	Email string    `json:"email"`
	OrgID uuid.UUID `json:"org_id"`
	Role  string    `json:"role"`
}
