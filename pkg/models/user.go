// Package models contains the structs used in the API, for both posting and returns.
package models

import (
	"github.com/google/uuid"
)

// return types

// UserID is a simple struct that contains a user's ID.
type UserID struct {
	UserID uuid.UUID `json:"user_id"`
}

// UserMetadata is all the information about a specific user.
type UserMetadata struct {
	UserID         uuid.UUID               `json:"user_id"`
	Email          string                  `json:"email"`
	EmailConfirmed bool                    `json:"email_confirmed"`
	HasPassword    bool                    `json:"has_password"`
	Username       *string                 `json:"username"`
	FirstName      *string                 `json:"first_name"`
	LastName       *string                 `json:"last_name"`
	PictureURL     *string                 `json:"picture_url"`
	Locked         bool                    `json:"locked"`
	Enabled        bool                    `json:"enabled"`
	MFAEnabled     bool                    `json:"mfa_enabled"`
	CanCreateOrgs  bool                    `json:"can_create_org"`
	CreatedAt      int64                   `json:"created_at"`
	LastActiveAt   int64                   `json:"last_active_at"`
	LegacyUserID   *string                 `json:"legacy_user_id"`
	OrgIDToOrgInfo *map[uuid.UUID]OrgInfo  `json:"org_id_to_org_info"`
	Metadata       *map[string]interface{} `json:"metadata"`
	Properties     *map[string]interface{} `json:"properties"`
}

// OrgInfo is the information about an organization a user is in.
type OrgInfo struct {
	OrgID            uuid.UUID        `json:"org_id"`
	OrgName          string           `json:"org_name"`
	OrgRoleStructure OrgRoleStructure `json:"org_role_structure"`
	UserRole         string           `json:"user_role"`
	AdditionalRoles  []string         `json:"additional_roles"`
	UserPermissions  []string         `json:"user_permissions"`
}

// UserList is a paged list of users. The actual fetched users are in the Users field, and the
// pagination information is in the other fields.
type UserList struct {
	TotalUsers     int            `json:"total_users"`
	CurrentPage    int            `json:"current_page"`
	PageSize       int            `json:"page_size"`
	HasMoreResults bool           `json:"has_more_results"`
	Users          []UserMetadata `json:"users"`
}

// post types

// CreateUserParams is the information needed to create a user. Email is required, but all other fields are optional.
// If password is missing the will be allowed to create one on their own. Username, Firstname, and Lastname are only
// used if you have the corresponding settings enabled in your dashboard.
type CreateUserParams struct {
	Email                          string                  `json:"email"`
	EmailConfirmed                 *bool                   `json:"email_confirmed,omitempty"`
	SendEmailToConfirmEmailAddress *bool                   `json:"send_email_to_confirm_email_address,omitempty"`
	Password                       *string                 `json:"password,omitempty"`
	Username                       *string                 `json:"username,omitempty"`
	FirstName                      *string                 `json:"first_name,omitempty"`
	LastName                       *string                 `json:"last_name,omitempty"`
	CanCreateOrgs                  *bool                   `json:"can_create_orgs,omitempty"`
	AskUserToUpdatePasswordOnLogin *bool                   `json:"ask_user_to_update_password_on_login,omitempty"`
	Metadata                       *map[string]interface{} `json:"metadata,omitempty"`
	Properties                     *map[string]interface{} `json:"properties,omitempty"`
	IgnoreDomainRestrictions       *bool                   `json:"ignore_domain_restrictions,omitempty"`
}

// MigrateUserParams is the information needed to migrate a user from another system. Email is required, but all other
// fields are optional. ExistingUserID will be saved in the LegacyUserID field in UserMetadata. If ExistingPasswordHash
// is provided, the user will be able to log in with their same password.
type MigrateUserParams struct {
	Email                          string                  `json:"email"`
	EmailConfirmed                 *bool                   `json:"email_confirmed,omitempty"`
	ExistingUserID                 *string                 `json:"existing_user_id,omitempty"`
	ExistingPasswordHash           *string                 `json:"existing_password_hash,omitempty"`
	ExistingMfaBase32EncodedSecret *string                 `json:"existing_mfa_base32_encoded_secret,omitempty"`
	Enabled                        *bool                   `json:"enabled,omitempty"`
	Username                       *string                 `json:"username,omitempty"`
	FirstName                      *string                 `json:"first_name,omitempty"`
	LastName                       *string                 `json:"last_name,omitempty"`
	PictureUrl                     *string                 `json:"picture_url,omitempty"`
	UpdatePasswordRequired         *bool                   `json:"update_password_required,omitempty"`
	Properties                     *map[string]interface{} `json:"properties,omitempty"`
}

// MigrateUserPasswordParams is the information needed to migrate a user's password from another system.
type MigrateUserPasswordParams struct {
	UserID       uuid.UUID `json:"user_id"`
	PasswordHash string    `json:"password_hash"`
}

// UpdateEmailParams is the information needed to update a user's email address.
type UpdateEmail struct {
	Email                    string `json:"new_email"`
	RequireEmailConfirmation bool   `json:"require_email_confirmation"`
}

// UpdateUserMetadata is the information needed to update a user's metadata. All fields are optional, we'll only update
// the ones that are provided.
type UpdateUserMetadata struct {
	Username               *string                 `json:"username,omitempty"`
	FirstName              *string                 `json:"first_name,omitempty"`
	LastName               *string                 `json:"last_name,omitempty"`
	PictureURL             *string                 `json:"picture_url,omitempty"`
	UpdatePasswordRequired *bool                   `json:"update_password_required,omitempty"`
	CanCreateOrgs          *bool                   `json:"can_create_orgs,omitempty"`
	Metadata               *map[string]interface{} `json:"metadata,omitempty"`
	Properties             *map[string]interface{} `json:"properties,omitempty"`
	LegacyUserID           *string                 `json:"legacy_user_id,omitempty"`
}

// UserQueryParams is the information needed to query a pageable list of users. If left blank, PageSize defaults to 10
// and PageNumber defaults to 0. EmailOrUsername is a dual-user field that will search for that string in either the
// email or username fields. IncludeOrgs defaults to false, but if set to true will include all orgs each user is in.
type UserQueryParams struct {
	PageSize        *int    `json:"page_size,omitempty"`
	PageNumber      *int    `json:"page_number,omitempty"`
	OrderBy         *string `json:"order_by,omitempty"`
	EmailOrUsername *string `json:"email_or_username,omitempty"`
	IncludeOrgs     *bool   `json:"include_orgs,omitempty"`
	LegacyUserID    *string `json:"legacy_user_id,omitempty"`
}

// UpdateUserPasswordParam is the information needed to update a user's password.
type UpdateUserPasswordParam struct {
	Password                       string `json:"password"`
	AskUserToUpdatePasswordOnLogin *bool  `json:"ask_user_to_update_password_on_login,omitempty"`
}
