package client

import (
	"fmt"
	"github.com/propelauth/propelauth-go/pkg/models"
	testHelpers "github.com/propelauth/propelauth-go/pkg/test"
	"testing"
)

func TestInitializations(t *testing.T) {
	// setup common test data

	_, public_key := testHelpers.GenerateRSAKeys()

	tokenVerificationMetadata := &models.TokenVerificationMetadata{
		VerifierKey: public_key,
		Issuer:      "issuertest",
	}

	// test initialization

	t.Run("test init with trailing slash fails", func(t *testing.T) {
		_, err := InitBaseAuth("https://auth.example.com/", "apikey", tokenVerificationMetadata)
		if err == nil {
			t.Errorf("NewClient should have returned an error, but did not")
		}
	})

	t.Run("test init with http and not https fails", func(t *testing.T) {
		_, err := InitBaseAuth("http://auth.example.com", "apikey", tokenVerificationMetadata)
		if err == nil {
			t.Errorf("NewClient should have returned an error about https, but did not")
		}
	})

}

func TestValidations(t *testing.T) {
	// setup common test data

	private_key, public_key := testHelpers.GenerateRSAKeys()

	userId := testHelpers.RandomUserID()
	org := testHelpers.RandomOrg("Admin", nil)
	org.UserInheritedRolesPlusCurrentRole = []string{"Admin", "Member"}
	org.UserPermissions = []string{"Read", "Write"}
	orgIdToOrgMemberInfo := testHelpers.OrgsToOrgIdMap([]models.OrgMemberInfoFromToken{org})
	user := models.UserFromToken{
		UserId:               userId,
		OrgIdToOrgMemberInfo: orgIdToOrgMemberInfo,
	}

	accessToken := testHelpers.CreateAccessToken(user, private_key)

	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	tokenVerificationMetadata := &models.TokenVerificationMetadata{
		VerifierKey: public_key,
		Issuer:      "issuertest",
	}

	client, err := InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadata)
	if err != nil {
		t.Errorf("NewClient returned an error, cannot even begin the tests: %s", err)
		return
	}

	// test valid access token with different user requirements

	t.Run("ValidateAccessTokenAndGetUser", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUser(authHeader)
		if err != nil {
			fmt.Println(authHeader)
			t.Errorf("ValidateAccessTokenAndGetUser returned an error: %s", err)
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrg", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUserWithOrg(authHeader, org.OrgId)
		if err != nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrg returned an error: %s", err)
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrgByMinimumRole", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUserWithOrgByMinimumRole(authHeader, org.OrgId, "Member")
		if err != nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByMinimumRole for Member returned an error: %s", err)
		}

		_, err = client.ValidateAccessTokenAndGetUserWithOrgByMinimumRole(authHeader, org.OrgId, "Owner")
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByMinimumRole for Owner should have returned an error.")
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrgByExactRole", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUserWithOrgByExactRole(authHeader, org.OrgId, "Admin")
		if err != nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByExactRole for Admin returned an error: %s", err)
		}

		_, err = client.ValidateAccessTokenAndGetUserWithOrgByExactRole(authHeader, org.OrgId, "Member")
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByExactRole for Member returned an error.")
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrgByPermission", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUserWithOrgByPermission(authHeader, org.OrgId, "Read")
		if err != nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByPermission for Read returned an error: %s", err)
		}

		_, err = client.ValidateAccessTokenAndGetUserWithOrgByPermission(authHeader, org.OrgId, "Delete")
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByPermission for Delete should have returned an error")
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrgByAllPermissions", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUserWithOrgByAllPermissions(authHeader, org.OrgId, []string{"Read", "Write"})
		if err != nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByAllPermissions for Read/Write returned an error: %s", err)
		}

		_, err = client.ValidateAccessTokenAndGetUserWithOrgByAllPermissions(authHeader, org.OrgId, []string{"Read", "Write", "Delete"})
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUserWithOrgByAllPermissions for Read/Write/Delete should have returned an error")
		}
	})

	// test bad headesr and bad access tokens

	t.Run("test basic validation fails Without Header", func(t *testing.T) {
		_, err := client.ValidateAccessTokenAndGetUser("")
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about the header")
		}
	})

	t.Run("test basic validation fails With Invalid Header", func(t *testing.T) {
		badAuthHeader := fmt.Sprintf("BadBearerHeader %s", accessToken)
		_, err := client.ValidateAccessTokenAndGetUser(badAuthHeader)
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about the header")
		}
	})

	t.Run("test basic validation fails With Wrong Token", func(t *testing.T) {
		badAuthHeader := "Bearer thisisafaketoken"
		_, err := client.ValidateAccessTokenAndGetUser(badAuthHeader)
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about the token")
		}
	})

	t.Run("test basic validation fails With Expired Token", func(t *testing.T) {
		// setup the expired token
		accessToken := testHelpers.CreateExpiredAccessToken(user, private_key)
		authHeader := fmt.Sprintf("Bearer %s", accessToken)

		// run the test
		_, err := client.ValidateAccessTokenAndGetUser(authHeader)
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about the token")
		}
	})

	t.Run("test basic validation fails With Bad Issuer", func(t *testing.T) {
		// setup the bad issuer
		tokenVerificationMetadata := &models.TokenVerificationMetadata{
			VerifierKey: public_key,
			Issuer:      "newissuertestthatwontmatch",
		}

		client, err := InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadata)
		if err != nil {
			t.Errorf("NewClient returned an error, cannot continue this test: %s", err)
			return
		}

		// run the test
		_, err = client.ValidateAccessTokenAndGetUser(authHeader)
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about issuer")
		}
	})

	t.Run("test basic validation fails With Wrong Key", func(t *testing.T) {
		// generate a client with new random keys
		_, public_key := testHelpers.GenerateRSAKeys()

		tokenVerificationMetadata := &models.TokenVerificationMetadata{
			VerifierKey: public_key,
			Issuer:      "issuertest",
		}

		client, err := InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadata)
		if err != nil {
			t.Errorf("NewClient returned an error, cannot even begin the tests: %s", err)
			return
		}

		// run the test
		_, err = client.ValidateAccessTokenAndGetUser(authHeader)
		if err == nil {
			t.Errorf("ValidateAccessTokenAndGetUser should have returned an error about decoding the token")
		}
	})
}
