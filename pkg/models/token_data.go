package models

import (
	"crypto/rsa"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// access token

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type AccessTokenResponse struct {
	AccessToken AccessTokenData `json:"access_token"`
}

type AccessTokenData struct {
	AccessToken          string                            `json:"access_token"`
	ExpiresAtSeconds     int64                             `json:"expires_at_seconds"`
	OrgIdToOrgMemberInfo map[string]OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
	User                 UserMetadata                      `json:"user"`
	ImpersonatorUser     UserID                            `json:"impersonator_user,omitempty"`
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

func (o *OrgMemberInfoFromToken) IsRole(exactRole string) bool {
	return exactRole == o.UserAssignedRole
}

func (o *OrgMemberInfoFromToken) IsAtLeastRole(minimumRoles string) bool {
	for _, role := range o.UserInheritedRolesPlusCurrentRole {
		if minimumRoles == role {
			return true
		}
	}
	return false
}

func (o *OrgMemberInfoFromToken) HasPermission(permission string) bool {
	for _, p := range o.UserPermissions {
		if permission == p {
			return true
		}
	}
	return false
}

// TODO - ridiculously not efficient
func (o *OrgMemberInfoFromToken) HasAllPermissions(permissions []string) bool {
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
