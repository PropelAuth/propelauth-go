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
	Username       string             `json:"username,omitempty"`
	FirstName      string             `json:"first_name,omitempty"`
	LastName       string             `json:"last_name,omitempty"`
	PictureURL     string             `json:"picture_url,omitempty"`
	Locked         bool               `json:"locked"`
	Enabled        bool               `json:"enabled"`
	MFAEnabled     bool               `json:"mfa_enabled"`
	CreatedAt      int64              `json:"created_at"`
	LastActiveAt   int64              `json:"last_active_at"`
	LegacyUserID   string             `json:"legacy_user_id,omitempty"`
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
	Email                          string `json:"email"`
	EmailConfirmed                 bool   `json:"email_confirmed"`
	SendEmailToConfirmEmailAddress bool   `json:"send_email_to_confirm_email_address"`
	Password                       string `json:"password"`
	Username                       string `json:"username"`
	FirstName                      string `json:"first_name"`
	LastName                       string `json:"last_name"`
}

type MigrateUserParams struct {
	Email                          string `json:"email"`
	EmailConfirmed                 bool   `json:"email_confirmed"`
	ExistingUserId                 string `json:"existing_user_id"`
	ExistingPasswordHash           string `json:"existing_password_hash"`
	ExistingMfaBase32EncodedSecret string `json:"existing_mfa_base32_encoded_secret"`
	Enabled                        bool   `json:"enabled"`
	Username                       string `json:"username"`
	FirstName                      string `json:"first_name"`
	LastName                       string `json:"last_name"`
}

type UpdateEmail struct {
	Email                    string `json:"email"`
	RequireEmailConfirmation bool   `json:"require_email_confirmation"`
}

type UpdateUserMetadata struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserQueryParams struct {
	PageSize        int    `json:"page_size"`
	PageNumber      int    `json:"page_number"`
	OrderBy         string `json:"order_by"`
	EmailOrUsername string `json:"email_or_username"`
	IncludeOrgs     bool   `json:"include_orgs"`
}

type UpdateUserPasswordParam struct {
	Password                       []string `json:"password"`
	AskUserToUpdatePasswordOnLogin []string `json:"askUserToUpdatePasswordOnLogin"`
}
