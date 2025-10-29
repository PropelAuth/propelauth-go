package models

import (
	"github.com/google/uuid"
)

// post types

// UserInOrgQueryParams is the information for querying a pageable list of users in an organization.
// If left blank, PageSize defaults to 10 and PageNumber defaults to 0. IncludeOrgs defaults to false,
// but if set to true will include all orgs each user is in. Role can be used to filter users by
// their role in the organization. Only case-sensitive, exact matches will be returned for Role.
type UserInOrgQueryParams struct {
	PageSize    *int    `json:"page_size,omitempty"`
	PageNumber  *int    `json:"page_number,omitempty"`
	IncludeOrgs *bool   `json:"include_orgs,omitempty"`
	Role        *string `json:"role,omitempty"`
}

// AddUserToOrg is the information needed to add a user to an organization. Role is just a string, but
// it has to match up to one of your defined roles, by default these are Owner, Admin, or Member, but
// they can be changed via your dashboard. If you've configured multiple roles per user in your project,
// you can also include additional_roles to add multiple roles to a user in an organization.
type AddUserToOrg struct {
	UserID          uuid.UUID `json:"user_id"`
	OrgID           uuid.UUID `json:"org_id"`
	Role            string    `json:"role"`
	AdditionalRoles []string  `json:"additional_roles,omitempty"`
}

// RemoveUserFromOrg is the information needed to remove a user from an organization.
type RemoveUserFromOrg struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
}

// InviteUserToOrg is the information needed to invite a new user to join an organization. Role is
// just a string, but it has to match up to one of your defined roles, by default these are Owner,
// Admin, or Member, but they can be changed via your dashboard. If you've configured multiple roles per user in your project,
// you can also include additional_roles to add multiple roles to a user in an organization.
type InviteUserToOrg struct {
	Email           string    `json:"email"`
	OrgID           uuid.UUID `json:"org_id"`
	Role            string    `json:"role"`
	AdditionalRoles []string  `json:"additional_roles,omitempty"`
}

type InviteUserToOrgByUserID struct {
	UserID          uuid.UUID `json:"user_id"`
	OrgID           uuid.UUID `json:"org_id"`
	Role            string    `json:"role"`
	AdditionalRoles []string  `json:"additional_roles,omitempty"`
}

// ChangeUserRoleInOrg is the information needed to change a user's role in an organization. Role is
// just a string, but it has to match up to one of your defined roles, by default these are Owner,
// Admin, or Member, but they can be changed via your dashboard. If you've configured multiple roles per user in your project,
// you can also include additional_roles to add multiple roles to a user in an organization.
type ChangeUserRoleInOrg struct {
	UserID          uuid.UUID `json:"user_id"`
	OrgID           uuid.UUID `json:"org_id"`
	Role            string    `json:"role"`
	AdditionalRoles []string  `json:"additional_roles,omitempty"`
}
