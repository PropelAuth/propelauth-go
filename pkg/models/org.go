package models

import (
	"github.com/google/uuid"
)

// return types

// OrgMetadata has the information about the organziation.
type OrgMetadata struct {
	OrgID            uuid.UUID              `json:"org_id"`
	Name             string                 `json:"org_name"`
	MaxUsers         *int                   `json:"max_users"`
	Metadata         map[string]interface{} `json:"metadata"`
	IsSamlConfigured bool                   `json:"is_saml_configured"`
}

// OrgList is a paged list of organizations. The actual fetched organizations are in the Orgs field, and the
// pagination information is in the other fields.
type OrgList struct {
	TotalOrgs      int           `json:"total_orgs"`
	CurrentPage    int           `json:"current_page"`
	PageSize       int           `json:"page_size"`
	HasMoreResults bool          `json:"has_more_results"`
	Orgs           []OrgMetadata `json:"orgs"`
}

// post types

// CreateOrg is the information needed to create an organization, which is just a name.
type CreateOrg struct {
	Name string `json:"name"`
}

// UpdateOrg is the information you can update in an organization. Each field is optional, we'll only update
// the fields you set. Note that AutojoinByDomain and RestrictToDomain require a validated domain, which can
// only be set in your dashboard.
type UpdateOrg struct {
	Name             *string                 `json:"name"`
	CanSetupSaml     *bool                   `json:"can_setup_saml"`
	AutojoinByDomain *bool                   `json:"autojoin_by_domain"`
	RestrictToDomain *bool                   `json:"restrict_to_domain"`
	MaxUsers         *int                    `json:"max_users"`
	Metadata         *map[string]interface{} `json:"metadata,omitempty"`
	Domain           *string                 `json:"domain,omitempty"`
	Require2FABy     *string                 `json:"require_2fa_by,omitempty"`
}

// OrgQueryParams is the information for querying a pageable organization list. If left blank, PageSize
// defaults to 10 and PageNumber defaults to 0.
type OrgQueryParams struct {
	PageSize   *int    `json:"page_size,omitempty"`
	PageNumber *int    `json:"page_number,omitempty"`
	OrderBy    *string `json:"order_by,omitempty"`
	Name       *string `json:"name,omitempty"`
}

// CreateOrgV2Params is the information needed to create an organization, as well as some optional fields.
type CreateOrgV2Params struct {
	Name                          string `json:"name"`
	Domain                        string `json:"domain,omitempty"`
	EnableAutoJoiningByDomain     bool   `json:"enable_auto_joining_by_domain,omitempty"`
	MembersMustHaveMatchingDomain bool   `json:"members_must_have_matching_domain,omitempty"`
	MaxUsers                      int    `json:"max_users,omitempty"`
}

// CreateOrgV2Response is the information returned when creating an organization.
type CreateOrgV2Response struct {
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
}
