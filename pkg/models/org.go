package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// return types

// OrgMetadata has the information about the organziation.
type OrgMetadata struct {
	OrgID                 uuid.UUID              `json:"org_id"`
	Name                  string                 `json:"org_name"`
	MaxUsers              *int                   `json:"max_users"`
	Metadata              map[string]interface{} `json:"metadata"`
	IsSamlConfigured      bool                   `json:"is_saml_configured"`
	CustomRoleMappingName *string                `json:"custom_role_mapping_name"`
	LegacyOrgId           string                 `json:"legacy_org_id"`
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

// CustomRoleMapping has the information about a Custom Role-to-Permissions.
type CustomRoleMapping struct {
	CustomRoleMappingName  string `json:"custom_role_mapping_name"`
	NumberOfOrgsSubscribed int    `json:"num_orgs_subscribed"`
}

// CustomRoleMappingList is a total list of all Custom Role-to-Permissions available
// in your environment.
type CustomRoleMappingList struct {
	CustomRoleMappings []CustomRoleMapping `json:"custom_role_mappings"`
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
	LegacyOrgId      *string                 `json:"legacy_org_id,omitempty"`
}

// OrgRoleMappingSubscription is the information needed to subscribe an organization to a
// Custom Role-to-Permissions mapping.
type OrgRoleMappingSubscription struct {
	CustomRoleMappingName string `json:"custom_role_mapping_name"`
}

// OrgQueryParams is the information for querying a pageable organization list. If left blank, PageSize
// defaults to 10 and PageNumber defaults to 0.
type OrgQueryParams struct {
	PageSize    *int    `json:"page_size,omitempty"`
	PageNumber  *int    `json:"page_number,omitempty"`
	OrderBy     *string `json:"order_by,omitempty"`
	Name        *string `json:"name,omitempty"`
	LegacyOrgId *string `json:"legacy_org_id,omitempty"`
}

// CreateOrgV2Params is the information needed to create an organization, as well as some optional fields.
type CreateOrgV2Params struct {
	Name                          string  `json:"name"`
	Domain                        string  `json:"domain,omitempty"`
	EnableAutoJoiningByDomain     bool    `json:"enable_auto_joining_by_domain,omitempty"`
	MembersMustHaveMatchingDomain bool    `json:"members_must_have_matching_domain,omitempty"`
	MaxUsers                      int     `json:"max_users,omitempty"`
	CustomRoleMappingName         *string `json:"custom_role_mapping_name,omitempty"`
	LegacyOrgId                   *string `json:"legacy_org_id,omitempty"`
}

// CreateOrgV2Response is the information returned when creating an organization.
type CreateOrgV2Response struct {
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
}

type OrgRoleStructure uint8

const (
	SingleRoleInHierarchy OrgRoleStructure = iota
	MultiRole
)

func (o OrgRoleStructure) String() string {
	return [...]string{"single_role_in_hierarchy", "multi_role"}[o]
}

func (o *OrgRoleStructure) FromString(s string) OrgRoleStructure {
	switch s {
	case "single_role_in_hierarchy":
		return SingleRoleInHierarchy
	case "multi_role":
		return MultiRole
	default:
		return SingleRoleInHierarchy
	}
}

func (o OrgRoleStructure) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *OrgRoleStructure) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*o = o.FromString(s)
	return nil
}

type PendingInvite struct {
	InviteeEmail         string     `json:"invitee_email"`
	OrgID                uuid.UUID  `json:"org_id"`
	OrgName              string     `json:"org_name"`
	RoleInOrg            string     `json:"role_in_org"`
	AdditionalRolesInOrg []string   `json:"additional_roles_in_org"`
	CreatedAt            int64      `json:"created_at"`
	ExpiresAt            int64      `json:"expires_at"`
	InviterEmail         *string    `json:"inviter_email"`
	InviterUserID        *uuid.UUID `json:"inviter_user_id"`
}

type PendingInvitesPage struct {
	TotalInvites   int             `json:"total_invites"`
	CurrentPage    int             `json:"current_page"`
	PageSize       int             `json:"page_size"`
	HasMoreResults bool            `json:"has_more_results"`
	Invites        []PendingInvite `json:"invites"`
}

type FetchPendingInvitesParams struct {
	PageSize   *int       `json:"page_size,omitempty"`
	PageNumber *int       `json:"page_number,omitempty"`
	OrgID      *uuid.UUID `json:"org_id,omitempty"`
}

type RevokePendingOrgInvite struct {
	OrgID        *uuid.UUID `json:"org_id,omitempty"`
	InviteeEmail string     `json:"invitee_email"`
}

// CreateSamlConnectionLinkBody is the information needed to create a SAML connection link.
type CreateSamlConnectionLinkBody struct {
	ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
}

type CreateSamlConnectionLinkResponse struct {
	URL string `json:"url"`
}
