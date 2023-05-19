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

type APIKeyValidation struct {
	Metadata      map[string]interface{} `json:"metadata"`
	UserMetadata  map[string]interface{} `json:"user_metadata"`
	OrgMetadata   map[string]interface{} `json:"org_metadata"`
	UserRoleInOrg string                 `json:"user_role_in_org"`
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
