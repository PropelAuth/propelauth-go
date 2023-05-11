package models

import (
	"crypto/rsa"
	"crypto/x509"
	pem "encoding/pem"
	"fmt"
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
	ImpersonatorUser     UserID                            `json:"impersonator_user,omitempty"`
}

// Models to hold public key data, that is used when initializing the client.

// TokenVerificationMetadataInput is a public key the user can pass in to initialize the client.
type TokenVerificationMetadataInput struct {
	VerifierKey string
	Issuer      string
}

// AuthTokenVerificationMetadataResponse is the response from the auth server when getting the public key.
type AuthTokenVerificationMetadataResponse struct {
	PublicKeyPem rsa.PublicKey `json:"public_key_pem"`
}

// TokenVerificationMetadata is the public key data needed to verify a token.
type TokenVerificationMetadata struct {
	VerifierKey rsa.PublicKey
	Issuer      string
}

// ConvertToTokenVerificationMetadata converts the public key from a string to a rsa.PublicKey, to make a TokenVerificationMetadata struct
func (o *TokenVerificationMetadataInput) ConvertToTokenVerificationMetadata() (*TokenVerificationMetadata, error) {
	pemBytes := []byte(o.VerifierKey)
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("Empty block found when decoding PEM block")
	} else if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("Wrong block type found when decoding PEM block, expecting PUBLIC KEY, found %s", block.Type)
	}

	parseResult, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Error when parsing public key: %w", err)
	}

	convertedToken := &TokenVerificationMetadata{
		VerifierKey: *parseResult,
		Issuer:      o.Issuer,
	}

	return convertedToken, nil
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
	LegacyUserID         string                             `json:"legacy_user_id,omitempty"`
	ImpersonatorUserID   uuid.UUID                          `json:"impersonator_user_id,omitempty"`
	Metadata             map[string]interface{}             `json:"metadata,omitempty"`
	OrgIDToOrgMemberInfo map[string]*OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
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
	UserAssignedRole                  string                 `json:"user_assigned_role"`
	UserInheritedRolesPlusCurrentRole []string               `json:"user_inherited_roles_plus_current_role"`
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
	// TODO - ridiculously not efficient
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
