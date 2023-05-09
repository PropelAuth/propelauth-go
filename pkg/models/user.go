package models

import (
	"github.com/google/uuid"
)

// return types

type UserID struct {
	UserID uuid.UUID `json:"user_id"`
}

type UserMetadata struct {
	UserID         uuid.UUID          `json:"user_id"`
	Email          string             `json:"email"`
	EmailConfirmed bool               `json:"email_confirmed"`
	HasPassword    bool               `json:"has_password"`
	Username       string             `json:"username,omempty"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	PictureURL     string             `json:"picture_url"`
	Locked         bool               `json:"locked"`
	Enabled        bool               `json:"enabled"`
	MFAEnabled     bool               `json:"mfa_enabled"`
	CreatedAt      int64              `json:"created_at"`
	LastActiveAt   int64              `json:"last_active_at"`
	LegacyUserID   string             `json:"legacy_user_id"`
	OrgIDToOrgInfo map[string]OrgInfo `json:"org_id_to_org_info"`
}

type OrgInfo struct {
	OrgID    uuid.UUID `json:"org_id"`
	OrgName  string    `json:"org_name"`
	UserRole string    `json:"user_role"`
}

type UserList struct {
	TotalUsers     int      `json:"total_users"`
	CurrentPage    int      `json:"current_page"`
	PageSize       int      `json:"page_size"`
	HasMoreResults bool     `json:"has_more_results"`
	Users          []UserID `json:"users"`
}

// post types

type CreateUserParams struct {
	Email                          string  `json:"email"`
	EmailConfirmed                 *bool   `json:"email_confirmed,omitempty"`
	SendEmailToConfirmEmailAddress *bool   `json:"send_email_to_confirm_email_address,omitempty"`
	Password                       *string `json:"password,omitempty"`
	Username                       *string `json:"username,omitempty"`
	FirstName                      *string `json:"first_name,omitempty"`
	LastName                       *string `json:"last_name,omitempty"`
}

type MigrateUserParams struct {
	Email                          string  `json:"email"`
	EmailConfirmed                 *bool   `json:"email_confirmed,omitempty"`
	ExistingUserId                 *string `json:"existing_user_id,omitempty"`
	ExistingPasswordHash           *string `json:"existing_password_hash,omitempty"`
	ExistingMfaBase32EncodedSecret *string `json:"existing_mfa_base32_encoded_secret,omitempty"`
	Enabled                        *bool   `json:"enabled,omitempty"`
	Username                       *string `json:"username,omitempty"`
	FirstName                      *string `json:"first_name,omitempty"`
	LastName                       *string `json:"last_name,omitempty"`
}

type UpdateEmail struct {
	Email                    string `json:"email"`
	RequireEmailConfirmation bool   `json:"require_email_confirmation"`
}

type UpdateUserMetadata struct {
	Username               *string `json:"username,omitempty"`
	FirstName              *string `json:"first_name,omitempty"`
	LastName               *string `json:"last_name,omitempty"`
	PictureUrl             *string `json:"picture_url,omitempty"`
	UpdatePasswordRequired *bool   `json:"update_password_required,omitempty"`
}

type UserQueryParams struct {
	PageSize        *int    `json:"page_size,omitempty"`
	PageNumber      *int    `json:"page_number,omitempty"`
	OrderBy         *string `json:"order_by,omitempty"`
	EmailOrUsername *string `json:"email_or_username,omitempty"`
	IncludeOrgs     *bool   `json:"include_orgs,omitempty"`
}

type UpdateUserPasswordParam struct {
	Password                       string `json:"password"`
	AskUserToUpdatePasswordOnLogin *bool  `json:"askUserToUpdatePasswordOnLogin,omitempty"`
}
