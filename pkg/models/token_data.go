package models

import (
	"crypto/rsa"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	access_token string `json:"access_token"`
}

// models used when initializing the client

type AuthTokenVerificationMetadataResponse struct {
	PublicKeyPem rsa.PublicKey `json:"public_key_pem"`
}

type TokenVerificationMetadata struct {
	VerifierKey rsa.PublicKey
	Issuer      string
}

// data from token

type UserAndOrgMemberInfoFromToken struct {
	User          UserFromToken
	OrgMemberInfo OrgMemberInfoFromToken
}

type UserFromToken struct {
	UserId               uuid.UUID                          `json:"user_id"`
	LegacyUserId         string                             `json:"legacy_user_id,omitempty"`
	ImpersonatorUserId   uuid.UUID                          `json:"impersonator_user_id,omitempty"`
	Metadata             map[string]interface{}             `json:"metadata,omitempty"`
	OrgIdToOrgMemberInfo map[string]*OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
	jwt.RegisteredClaims
}

func (o *UserFromToken) GetOrgMemberInfo(orgId uuid.UUID) *OrgMemberInfoFromToken {
	return o.OrgIdToOrgMemberInfo[orgId.String()]
}

type OrgMemberInfoFromToken struct {
	OrgId                             uuid.UUID              `json:"org_id"`
	OrgName                           string                 `json:"org_name"`
	OrgMetadata                       map[string]interface{} `json:"org_metadata,omitempty"`
	UserAssignedRole                  string                 `json:"user_assigned_role"`
	UserInheritedRolesPlusCurrentRole []string               `json:"user_inherited_roles_plus_current_role"`
	UserPermissions                   []string               `json:"user_permissions,omitempty"`
}

func (o *OrgMemberInfoFromToken) VerifyExactRole(exactRole string) bool {
	if exactRole == o.UserAssignedRole {
		return true
	}
	return false
}

func (o *OrgMemberInfoFromToken) VerifyMinimumRole(minimumRoles string) bool {
	for _, role := range o.UserInheritedRolesPlusCurrentRole {
		if minimumRoles == role {
			return true
		}
	}
	return false
}

func (o *OrgMemberInfoFromToken) VerifyPermission(permission string) bool {
	for _, p := range o.UserPermissions {
		if permission == p {
			return true
		}
	}
	return false
}

// TODO - ridiculously not efficient
func (o *OrgMemberInfoFromToken) VerifyAllPermissions(permissions []string) bool {
	for _, permission := range permissions {
		found := false
		for _, p := range o.UserPermissions {
			if permission == p {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
