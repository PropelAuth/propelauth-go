package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	pem "encoding/pem"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/propelauth/propelauth-go/pkg/models"
)

func RandomUserID() uuid.UUID {
	return uuid.New()
}

func RandomOrgID() uuid.UUID {
	return uuid.New()
}

// Represents the incoming JSON from the auth server.
func RandomOrg(userRole string, multi_role bool) models.OrgMemberInfoFromToken {
	if multi_role {
		return models.OrgMemberInfoFromToken{
			OrgID:                             RandomOrgID(),
			OrgName:                           "orgname",
			OrgMetadata:                       map[string]interface{}{},
			OrgRoleStructure:                  models.MultiRole,
			UserAssignedRole:                  userRole,
			UserInheritedRolesPlusCurrentRole: []string{userRole},
		}
	} else {
		return models.OrgMemberInfoFromToken{
			OrgID:                             RandomOrgID(),
			OrgName:                           "orgname",
			OrgMetadata:                       map[string]interface{}{},
			UserAssignedRole:                  userRole,
			UserInheritedRolesPlusCurrentRole: []string{userRole},
		}
	}
}

// Convert a list of orgs to a map of org_id to org, which is used in UserFromToken.
func OrgsToOrgIDMap(orgs []models.OrgMemberInfoFromToken) map[string]*models.OrgMemberInfoFromToken {
	orgIDToOrgMemberInfo := make(map[string]*models.OrgMemberInfoFromToken)

	for i := range orgs {
		uuidAsString := orgs[i].OrgID.String()
		orgIDToOrgMemberInfo[uuidAsString] = &orgs[i]
	}

	return orgIDToOrgMemberInfo
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
func GenerateRSAKeys() (*rsa.PrivateKey, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// convert privateKey.PublicKey to a string, which is what our users are likely to have
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return privateKey, string(publicKeyPem)
}
