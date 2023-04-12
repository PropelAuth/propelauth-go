package client

import (
	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID                `json:"user_id"`
	LegacyUserID   string                   `json:"legacy_user_id"`
	OrgIDToOrgInfo map[string]UserInOrgInfo `json:"org_id_to_org_info"`
}

type UserInOrgInfo struct {
	OrgID                             uuid.UUID `json:"org_id"`
	OrgName                           string    `json:"org_name"`
	UserAssignedRole                  string    `json:"user_assigned_role"`
	UserInheritedRolesPlusCurrentRole []string  `json:"user_inherited_roles_plus_current_role"`
	UserPermissions                   []string  `json:"user_permissions"`
}
