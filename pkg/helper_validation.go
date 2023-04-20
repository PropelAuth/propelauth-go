package client

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
)

type ValidationHelper struct{}

func (o *ValidationHelper) ValidateAccessTokenAndGetUser(accessToken string, tokenVerificationMetadata TokenVerificationMetadata) (*UserFromToken, error) {

	claims := &UserFromToken{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Error decoding JWT: Unexpected signing method: %v", token.Header["alg"])
		}

		//verifierBytes := []byte(tokenVerificationMetadata.VerifierKey)
		return tokenVerificationMetadata.VerifierKey, nil
	})

	// friendly error messages
	if err == jwt.ErrTokenMalformed {
		return nil, fmt.Errorf("Error decoding JWT: malformed token")
	} else if err == jwt.ErrTokenSignatureInvalid {
		return nil, fmt.Errorf("Error decoding JWT: invalid token signature")
	} else if err == jwt.ErrTokenExpired || err == jwt.ErrTokenNotValidYet {
		return nil, fmt.Errorf("Error decoding JWT: expired token")
	} else if err != nil {
		return nil, fmt.Errorf("Error decoding JWT: unknown error: %v", err)
	} else if !token.Valid {
		return nil, fmt.Errorf("Error decoding JWT: invalid token")
	} else if claims.Issuer != tokenVerificationMetadata.Issuer {
		return nil, fmt.Errorf("Error decoding JWT: invalid issuer")
	}

	return claims, nil
}

func (o *ValidationHelper) ExtractTokenFromAuthorizationHeader(authHeader string) (string, error) {
	split := strings.Split(authHeader, " ")

	if len(split) != 2 {
		return "", fmt.Errorf("Authorization header is not in the correct format")
	}
	if split[0] != "Bearer" {
		return "", fmt.Errorf("Authorization header is not in the correct format")
	}

	return split[1], nil
}

func (o *ValidationHelper) ValidateOrgAccessAndGetOrgMemberInfo(user *UserFromToken, orgId uuid.UUID) (*OrgMemberInfoFromToken, error) {
	if user.OrgIdToOrgMemberInfo == nil {
		return nil, fmt.Errorf("User does not have access to any organizations")
	}

	orgMemberInfo := user.GetOrgMemberInfo(orgId)

	if orgMemberInfo == nil {
		return nil, fmt.Errorf("User does not have access to this organization")
	}

	return orgMemberInfo, nil
}

func (o *ValidationHelper) ValidateOrgAccessAndGetOrgMemberInfoByMinimumRole(user *UserFromToken, orgId uuid.UUID, minimumRole string) (*OrgMemberInfoFromToken, error) {
	orgMemberInfo, err := o.ValidateOrgAccessAndGetOrgMemberInfo(user, orgId)
	if err != nil {
		return nil, err
	}

	if !orgMemberInfo.VerifyMinimumRole(minimumRole) {
		return nil, fmt.Errorf("User does not have minimum role needed in this organization")
	}

	return orgMemberInfo, nil
}

func (o *ValidationHelper) ValidateOrgAccessAndGetOrgMemberInfoByExactRole(user *UserFromToken, orgId uuid.UUID, exactRole string) (*OrgMemberInfoFromToken, error) {
	orgMemberInfo, err := o.ValidateOrgAccessAndGetOrgMemberInfo(user, orgId)
	if err != nil {
		return nil, err
	}

	if !orgMemberInfo.VerifyExactRole(exactRole) {
		return nil, fmt.Errorf("User does not have the exact role needed in this organization")
	}

	return orgMemberInfo, nil
}

func (o *ValidationHelper) ValidateOrgAccessAndGetOrgMemberInfoByPermission(user *UserFromToken, orgId uuid.UUID, permission string) (*OrgMemberInfoFromToken, error) {
	orgMemberInfo, err := o.ValidateOrgAccessAndGetOrgMemberInfo(user, orgId)
	if err != nil {
		return nil, err
	}

	if !orgMemberInfo.VerifyPermission(permission) {
		return nil, fmt.Errorf("User does not have the permission needed in this organization")
	}

	return orgMemberInfo, nil
}

func (o *ValidationHelper) ValidateOrgAccessAndGetOrgMemberInfoByAllPermissions(user *UserFromToken, orgId uuid.UUID, permissions []string) (*OrgMemberInfoFromToken, error) {
	orgMemberInfo, err := o.ValidateOrgAccessAndGetOrgMemberInfo(user, orgId)
	if err != nil {
		return nil, err
	}

	if !orgMemberInfo.VerifyAllPermissions(permissions) {
		return nil, fmt.Errorf("User does not have all the permissions needed in this organization")
	}

	return orgMemberInfo, nil
}
