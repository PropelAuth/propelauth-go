package models

import (
	"crypto/rsa"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Access token models.
// These are used behind the scenes in the client, and you probably won't need to use them directly.

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type AccessTokenResponse struct {
	AccessToken AccessTokenData `json:"access_token"`
}

type AccessTokenData struct {
	AccessToken          string                            `json:"access_token"`
	ExpiresAtSeconds     int64                             `json:"expires_at_seconds"`
	OrgIDToOrgMemberInfo map[string]OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
	User                 UserMetadata                      `json:"user"`
	ImpersonatorUser     *UserID                           `json:"impersonator_user,omitempty"`
}

// Models to hold public key data, that is used when initializing the client.

// TokenVerificationMetadataInput is a public key type the user can pass in to initialize the client. The public key is a string.
type TokenVerificationMetadataInput struct {
	VerifierKey string
	Issuer      string
}

// AuthTokenVerificationMetadataResponse is the response from the auth server when getting the public key.
type AuthTokenVerificationMetadataResponse struct {
	VerifierKeyPem string `json:"verifier_key_pem"`
}

// TokenVerificationMetadata is the public key type we use internally.
type TokenVerificationMetadata struct {
	VerifierKey rsa.PublicKey
	Issuer      string
}

// Data from token

// UserAndOrgMemberInfoFromToken is the user and organization data from the JWT.
type UserAndOrgMemberInfoFromToken struct {
	User          UserFromToken
	OrgMemberInfo OrgMemberInfoFromToken
}

// UserFromToken is the user data from the JWT.
type UserFromToken struct {
	UserID               uuid.UUID                          `json:"user_id"`
	LegacyUserID         *string                            `json:"legacy_user_id,omitempty"`
	ImpersonatorUserID   *uuid.UUID                         `json:"impersonator_user_id,omitempty"`
	OrgIDToOrgMemberInfo map[string]*OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
	Email                *string                            `json:"email"`
	FirstName            *string                            `json:"first_name,omitempty"`
	LastName             *string                            `json:"last_name,omitempty"`
	Username             *string                            `json:"username,omitempty"`
	Properties           map[string]interface{}             `json:"properties,omitempty"`
	jwt.RegisteredClaims
}

// GetOrgMemberInfo returns the OrgMemberInfoFromToken for the given Organization UUID.
func (o *UserFromToken) GetOrgMemberInfo(orgID uuid.UUID) *OrgMemberInfoFromToken {
	return o.OrgIDToOrgMemberInfo[orgID.String()]
}

// OrgMemberInfoFromToken is data about an organization and about this user's membership in it.
type OrgMemberInfoFromToken struct {
	OrgID                             uuid.UUID              `json:"org_id"`
	OrgName                           string                 `json:"org_name"`
	OrgMetadata                       map[string]interface{} `json:"org_metadata,omitempty"`
	URLSafeOrgName                    string                 `json:"url_safe_org_name,omitempty"`
	UserAssignedRole                  string                 `json:"user_role"`
	UserInheritedRolesPlusCurrentRole []string               `json:"inherited_user_roles_plus_current_role"`
	UserPermissions                   []string               `json:"user_permissions,omitempty"`
}

// IsRole returns true if the user has the exact role.
func (o *OrgMemberInfoFromToken) IsRole(exactRole string) bool {
	return exactRole == o.UserAssignedRole
}

// IsAtLeastRole returns true if the user has the exact role or a role that is higher in the hierarchy.
func (o *OrgMemberInfoFromToken) IsAtLeastRole(minimumRoles string) bool {
	for _, role := range o.UserInheritedRolesPlusCurrentRole {
		if minimumRoles == role {
			return true
		}
	}

	return false
}

// HasPermission returns true if the user has the exact permission.
func (o *OrgMemberInfoFromToken) HasPermission(permission string) bool {
	for _, p := range o.UserPermissions {
		if permission == p {
			return true
		}
	}

	return false
}

// HasAllPermissions returns true if the user has all of the permissions.
func (o *OrgMemberInfoFromToken) HasAllPermissions(permissions []string) bool {
	permissionsSet := make(map[string]bool)
	for _, p := range o.UserPermissions {
		permissionsSet[p] = true
	}

	for _, permission := range permissions {
		if !permissionsSet[permission] {
			return false
		}
	}

	return true
}
