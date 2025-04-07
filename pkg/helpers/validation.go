// Package helpers contains internal helper functions for the client, and are not intended to be used directly.
package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	pem "encoding/pem"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	models "github.com/propelauth/propelauth-go/pkg/models"
)

type ValidationHelperInterface interface {
	ValidateAccessTokenAndGetUser(accessToken string, tokenVerificationMetadata models.TokenVerificationMetadata) (*models.UserFromToken, error)
	ExtractTokenFromAuthorizationHeader(authHeader string) (string, error)
	ConvertPEMStringToRSAPublicKey(pemString string) (*rsa.PublicKey, error)
}

type ValidationHelper struct{}

// ValidateAccessTokenAndGetUser validates the access token and returns the user data. Instead of using this
// directly, look at client.GetUser(authHeader) instead.
func (o *ValidationHelper) ValidateAccessTokenAndGetUser(accessToken string, tokenVerificationMetadata models.TokenVerificationMetadata) (*models.UserFromToken, error) {
	userFromToken := &models.UserFromToken{}

	token, err := jwt.ParseWithClaims(accessToken, userFromToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Error decoding JWT: Unexpected signing method: %v", token.Header["alg"])
		}

		return &tokenVerificationMetadata.VerifierKey, nil
	},
		jwt.WithIssuer(tokenVerificationMetadata.Issuer),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	)

	// friendly error messages
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("Error decoding JWT: malformed token")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, fmt.Errorf("Error decoding JWT: invalid token signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, fmt.Errorf("Error decoding JWT: expired token")
	} else if err != nil {
		return nil, fmt.Errorf("Error decoding JWT: unknown error: %w", err)
	} else if !token.Valid {
		return nil, fmt.Errorf("Error decoding JWT: invalid token")
	} else if userFromToken.Issuer != tokenVerificationMetadata.Issuer {
		return nil, fmt.Errorf("Error decoding JWT: invalid issuer")
	}

	userFromTokenWithActiveOrg := AssignActiveOrg(userFromToken)

	// if the login method is nil, set it to unknown
	if userFromTokenWithActiveOrg.LoginMethod == nil {
		userFromTokenWithActiveOrg.LoginMethod = &models.LoginMethod{
			LoginMethod: "unknown",
		}
	}

	return userFromTokenWithActiveOrg, nil
}

func AssignActiveOrg(userFromToken *models.UserFromToken) *models.UserFromToken {
	// Properly assign OrgIdToOrgMemberInfo for Active Org Support
	if userFromToken.OrgMemberInfo != nil {
		userFromToken.ActiveOrgId = &userFromToken.OrgMemberInfo.OrgID
		userFromToken.OrgIDToOrgMemberInfo = map[string]*models.OrgMemberInfoFromToken{
			userFromToken.OrgMemberInfo.OrgID.String(): userFromToken.OrgMemberInfo,
		}
		userFromToken.OrgMemberInfo = nil
	}
	return userFromToken

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

// ConvertToTokenVerificationMetadata converts the public key from a string to a rsa.PublicKey, to make a TokenVerificationMetadata struct
func (o *ValidationHelper) ConvertPEMStringToRSAPublicKey(pemString string) (*rsa.PublicKey, error) {
	pemBytes := []byte(pemString)
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("Empty block found when decoding PEM block")
	} else if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("Wrong block type found when decoding PEM block, expecting PUBLIC KEY, found %s", block.Type)
	}

	// TODO: There's got to be a better way to do this

	// first try PKCS1
	parseResult, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return parseResult, nil
	}

	// then try PKIX
	parseResult2, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		return parseResult2.(*rsa.PublicKey), nil
	}

	return nil, fmt.Errorf("Error when parsing public key: %w", err)
}
