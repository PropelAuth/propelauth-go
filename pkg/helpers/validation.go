// Package helpers contains internal helper functions for the client, and are not intended to be used directly.
package helpers

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/propelauth/propelauth-go/pkg/models"
	"strings"
)

type ValidationHelperInterface interface {
	ValidateAccessTokenAndGetUser(accessToken string, tokenVerificationMetadata models.TokenVerificationMetadata) (*models.UserFromToken, error)
	ExtractTokenFromAuthorizationHeader(authHeader string) (string, error)
}

type ValidationHelper struct{}

// ValidateAccessTokenAndGetUser validates the access token and returns the user data. Instead of using this
// directly, look at client.GetUser(authHeader) instead.
func (o *ValidationHelper) ValidateAccessTokenAndGetUser(accessToken string, tokenVerificationMetadata models.TokenVerificationMetadata) (*models.UserFromToken, error) {

	claims := &models.UserFromToken{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Error decoding JWT: Unexpected signing method: %v", token.Header["alg"])
		}

		return &tokenVerificationMetadata.VerifierKey, nil
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
