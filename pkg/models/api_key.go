package models

import (
	"github.com/google/uuid"
)

// return types

type APIKeyFull struct {
	APIKeyId         string                 `json:"api_key_id"`
	CreatedAt        int                    `json:"created_at"`
	ExpiresAtSeconds int                    `json:"expires_at_seconds"`
	Metadata         map[string]interface{} `json:"metadata"`
	UserID           uuid.UUID              `json:"user_id"`
	OrgID            uuid.UUID              `json:"org_id"`
}

type APIKeyResultPage struct {
	APIKeys        []APIKeyFull `json:"api_keys"`
	TotalAPIKeys   int          `json:"total_api_keys"`
	CurrentPage    int          `json:"current_page"`
	PageSize       int          `json:"page_size"`
	HasMoreResults bool         `json:"has_more_results"`
}

type APIKeyOrgMetadata struct {
	OrgID        uuid.UUID              `json:"org_id"`
	OrgName      string                 `json:"org_name"`
	CanSetupSaml bool                   `json:"can_setup_saml"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type APIKeyValidation struct {
	Metadata  map[string]interface{}  `json:"metadata"`
	User      *UserMetadata           `json:"user"`
	Org       *APIKeyOrgMetadata      `json:"org"`
	UserInOrg *OrgMemberInfoFromToken `json:"user_in_org"`
}

type PersonalAPIKeyValidation struct {
	Metadata map[string]interface{} `json:"metadata"`
	User     UserMetadata           `json:"user"`
}

type OrgAPIKeyValidation struct {
	Metadata  map[string]interface{}  `json:"metadata"`
	Org       APIKeyOrgMetadata       `json:"org"`
	User      *UserMetadata           `json:"user"`
	UserInOrg *OrgMemberInfoFromToken `json:"user_in_org"`
}

type APIKeyNew struct {
	APIKeyID    string `json:"api_key_id"`
	APIKeyToken string `json:"api_key_token"`
}

// post types

type APIKeysQueryParams struct {
	OrgID      *uuid.UUID `json:"org_id,omitempty"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	UserEmail  *string    `json:"user_email,omitempty"`
	PageSize   *int       `json:"page_size,omitempty"`
	PageNumber *int       `json:"page_number,omitempty"`
}

type APIKeyCreateParams struct {
	OrgID            *uuid.UUID              `json:"org_id,omitempty"`
	UserID           *uuid.UUID              `json:"user_id,omitempty"`
	ExpiresAtSeconds *int                    `json:"expires_at_seconds,omitempty"`
	Metadata         *map[string]interface{} `json:"metadata,omitempty"`
}

type APIKeyUpdateParams struct {
	ExpiresAtSeconds *int                    `json:"expires_at_seconds,omitempty"`
	Metadata         *map[string]interface{} `json:"metadata,omitempty"`
}

type ApiKeyRateLimitError struct {
    WaitSeconds float64 `json:"wait_seconds"`
	ErrorCode   string  `json:"error_code"`
	UserFacingError string `json:"user_facing_error"`
}

func (e *ApiKeyRateLimitError) Error() string {
    return e.UserFacingError
}