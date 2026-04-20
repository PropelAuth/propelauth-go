package models

import (
	"github.com/google/uuid"
)

// return types

type FetchOrgScimGroupsRequest struct {
	OrgID      uuid.UUID  `json:"org_id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	PageSize   *int       `json:"page_size,omitempty"`
	PageNumber *int       `json:"page_number,omitempty"`
}

type FetchScimGroupRequest struct {
	OrgID             uuid.UUID `json:"org_id"`
	GroupID           uuid.UUID `json:"group_id"`
	MembersPageSize   *int      `json:"members_page_size,omitempty"`
	MembersPageNumber *int      `json:"members_page_number,omitempty"`
}

type ScimGroupResult struct {
	GroupID           uuid.UUID `json:"group_id"`
	DisplayName       string    `json:"display_name"`
	ExternalIDFromIDP string    `json:"external_id_from_idp"`
}

type ScimGroupResultPage struct {
	Groups      []ScimGroupResult `json:"groups"`
	TotalGroups int               `json:"total_groups"`
	PageNumber  int               `json:"page_number"`
	PageSize    int               `json:"page_size"`
}

type ScimGroupMember struct {
	UserID uuid.UUID `json:"user_id"`
}

type ScimGroup struct {
	GroupID           uuid.UUID         `json:"group_id"`
	DisplayName       string            `json:"display_name"`
	ExternalIDFromIDP string            `json:"external_id_from_idp"`
	Members           []ScimGroupMember `json:"members"`
}
