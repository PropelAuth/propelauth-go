package test

import (
	"crypto/rand"
	"crypto/rsa"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/propelauth/propelauth-go/pkg/models"
	"time"
)

func RandomUserID() uuid.UUID {
	return uuid.New()
}

func RandomOrgID() uuid.UUID {
	return uuid.New()
}

// Represents the incoming JSON from the auth server.
func RandomOrg(userRole string, permissions []string) models.OrgMemberInfoFromToken {
	return models.OrgMemberInfoFromToken{
		OrgId:                             RandomOrgID(),
		OrgName:                           "orgname",
		OrgMetadata:                       map[string]interface{}{},
		UserAssignedRole:                  userRole,
		UserInheritedRolesPlusCurrentRole: []string{userRole},
		UserPermissions:                   permissions,
	}
}

// Convert a list of orgs to a map of org_id to org, which is used in UserFromToken.
func OrgsToOrgIdMap(orgs []models.OrgMemberInfoFromToken) map[string]*models.OrgMemberInfoFromToken {
	orgIdToOrgMemberInfo := make(map[string]*models.OrgMemberInfoFromToken)
	for _, org := range orgs {
		uuidAsString := org.OrgId.String()
		orgIdToOrgMemberInfo[uuidAsString] = &org
	}
	return orgIdToOrgMemberInfo
}

// Create a JWT access token with the UserFromToken data.
func CreateAccessToken(user models.UserFromToken, privateKeyPem *rsa.PrivateKey) string {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "issuertest",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	tokenString, err := token.SignedString(privateKeyPem)
	if err != nil {
		panic("createAccessToken: error signing token: " + err.Error())
	}

	return tokenString
}

// Create an expired JWT access token with the UserFromToken data.
func CreateExpiredAccessToken(user models.UserFromToken, privateKeyPem *rsa.PrivateKey) string {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		Issuer:    "issuertest",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	tokenString, err := token.SignedString(privateKeyPem)
	if err != nil {
		panic("createExpiredAccessToken: error signing token: " + err.Error())
	}

	return tokenString
}

// Generate an RSA key pair.
func GenerateRSAKeys() (*rsa.PrivateKey, rsa.PublicKey) {
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return private_key, private_key.PublicKey
}
