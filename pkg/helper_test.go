package client

import (
	"crypto/rand"
	"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func randomUserID() uuid.UUID {
	return uuid.New()
}

func randomOrgID() uuid.UUID {
	return uuid.New()
}

// Represents the incoming JSON from the auth server.
func randomOrg(userRole string, permissions []string) OrgMemberInfoFromToken {
	return OrgMemberInfoFromToken{
		OrgId:                             randomOrgID(),
		OrgName:                           "orgname",
		OrgMetadata:                       map[string]interface{}{},
		UserAssignedRole:                  userRole,
		UserInheritedRolesPlusCurrentRole: []string{userRole},
		UserPermissions:                   permissions,
	}
}

// Convert alist of orgs to a map of org_id to org, which is used in UserFromToken.
func orgsToOrgIdMap(orgs []OrgMemberInfoFromToken) map[string]*OrgMemberInfoFromToken {
	orgIdToOrgMemberInfo := make(map[string]*OrgMemberInfoFromToken)
	for _, org := range orgs {
		uuidAsString := org.OrgId.String()
		orgIdToOrgMemberInfo[uuidAsString] = &org
	}
	return orgIdToOrgMemberInfo
}

// Create a JWT access token with the UserFromToken data.
func createAccessToken(user UserFromToken, privateKeyPem *rsa.PrivateKey) string {
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

	fmt.Println("createAccessToken result: " + tokenString)
	return tokenString
}

// Create a JWT access token with the UserFromToken data.
func createExpiredAccessToken(user UserFromToken, privateKeyPem *rsa.PrivateKey) string {
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

	fmt.Println("createExpiredAccessToken result: " + tokenString)
	return tokenString
}

// Generate an RSA key pair.
func generateRSAKeys() (*rsa.PrivateKey, rsa.PublicKey) {
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return private_key, private_key.PublicKey
}
