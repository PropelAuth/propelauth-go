package client_test

import (
	"fmt"
	"testing"

	propelauth "github.com/propelauth/propelauth-go/pkg"
	"github.com/propelauth/propelauth-go/pkg/models"
	testHelpers "github.com/propelauth/propelauth-go/pkg/test"
)

func TestInitializations(t *testing.T) {
	// setup common test data

	_, publicKey := testHelpers.GenerateRSAKeys()

	tokenVerificationMetadataInput := &models.TokenVerificationMetadataInput{
		VerifierKey: publicKey,
		Issuer:      "issuertest",
	}

	// test initialization

	t.Run("test init with trailing slash fails", func(t *testing.T) {
		_, err := propelauth.InitBaseAuth("https://auth.example.com/", "apikey", tokenVerificationMetadataInput)
		if err == nil {
			t.Errorf("NewClient should have returned an error, but did not")
		}
	})

	t.Run("test init with http and not https fails", func(t *testing.T) {
		_, err := propelauth.InitBaseAuth("http://auth.example.com", "apikey", tokenVerificationMetadataInput)
		if err == nil {
			t.Errorf("NewClient should have returned an error about https, but did not")
		}
	})
}

func TestValidations(t *testing.T) {
	// setup common test data

	privateKey, publicKey := testHelpers.GenerateRSAKeys()

	userID := testHelpers.RandomUserID()

	orgSingleRole := testHelpers.RandomOrg("Admin", false)
	orgSingleRole.UserInheritedRolesPlusCurrentRole = []string{"Admin", "Member"}
	orgSingleRole.UserPermissions = []string{"Read", "Write"}

	orgMultiRole := testHelpers.RandomOrg("Role A", true)
	orgMultiRole.UserAssignedAdditionalRoles = []string{"Role B", "Role C"}

	orgIDToOrgMemberInfo := testHelpers.OrgsToOrgIDMap([]models.OrgMemberInfoFromToken{orgSingleRole, orgMultiRole})

	user := models.UserFromToken{
		UserID:               userID,
		OrgIDToOrgMemberInfo: orgIDToOrgMemberInfo,
	}

	accessToken := testHelpers.CreateAccessToken(user, privateKey)

	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	tokenVerificationMetadataInput := &models.TokenVerificationMetadataInput{
		VerifierKey: publicKey,
		Issuer:      "issuertest",
	}

	client, err := propelauth.InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadataInput)
	if err != nil {
		t.Errorf("NewClient returned an error, cannot even begin the tests: %s", err)

		return
	}

	// test valid access token with different user requirements

	t.Run("GetUser", func(t *testing.T) {
		_, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser on %s returned an error: %s", authHeader, err)
		}
	})

	t.Run("ValidateAccessTokenAndGetUserWithOrg", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		// run tests

		orgMemberInfo := user.GetOrgMemberInfo(orgSingleRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}
	})

	t.Run("IsAtLeastRole, single_role_in_hierarchy", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgSingleRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.IsAtLeastRole("Member")
		if !result {
			t.Errorf("IsAtLeastRole with Member should have returned true")
		}

		result = orgMemberInfo.IsAtLeastRole("Admin")
		if !result {
			t.Errorf("IsAtLeastRole with Admin should have returned true")
		}

		result = orgMemberInfo.IsAtLeastRole("Owner")
		if result {
			t.Errorf("IsAtLeastRole with Owner should have returned false")
		}
	})

	t.Run("IsAtLeastRole, multi_role", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgMultiRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.IsAtLeastRole("Role A")
		if !result {
			t.Errorf("IsAtLeastRole with Role A should have returned true")
		}

		result = orgMemberInfo.IsAtLeastRole("Role C")
		if !result {
			t.Errorf("IsAtLeastRole with Role C should have returned true")
		}

		result = orgMemberInfo.IsAtLeastRole("Role D")
		if result {
			t.Errorf("IsAtLeastRole with Role D should have returned false")
		}
	})

	t.Run("IsRole, multi_role", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgMultiRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.IsRole("Role D")
		if result {
			t.Errorf("IsRole with Role D should have returned false")
		}

		result = orgMemberInfo.IsRole("Role A")
		if !result {
			t.Errorf("IsRole with Role A should have returned true")
		}

		result = orgMemberInfo.IsRole("Role B")
		if !result {
			t.Errorf("IsRole with Role B should have returned true")
		}
	})

	t.Run("IsRole, single_role_in_hierarchy", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgSingleRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.IsRole("Member")
		if result {
			t.Errorf("IsRole with Member should have returned false")
		}

		result = orgMemberInfo.IsRole("Admin")
		if !result {
			t.Errorf("IsRole with Admin should have returned true")
		}

		result = orgMemberInfo.IsRole("Owner")
		if result {
			t.Errorf("IsRole with Owner should have returned false")
		}
	})

	t.Run("HasPermission", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgSingleRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.HasPermission("Read")
		if !result {
			t.Errorf("HasPermission with Read should have returned true")
		}

		result = orgMemberInfo.HasPermission("Edit")
		if result {
			t.Errorf("HasPermission with Edit should have returned false")
		}
	})

	t.Run("HasAllPermissions", func(t *testing.T) {
		// setup tests

		user, err := client.GetUser(authHeader)
		if err != nil {
			t.Errorf("GetUser returned an error: %s", err)
		}

		orgMemberInfo := user.GetOrgMemberInfo(orgSingleRole.OrgID)
		if orgMemberInfo == nil {
			t.Errorf("GetOrgMemberInfo should have returned something")
		}

		// run tests

		result := orgMemberInfo.HasAllPermissions([]string{"Read", "Write"})
		if !result {
			t.Errorf("HasPermission with Read/Write should have returned true")
		}

		result = orgMemberInfo.HasAllPermissions([]string{"Read", "Write", "Delete"})
		if result {
			t.Errorf("HasPermission with Read/Write/Delete should have returned false")
		}
	})

	// test bad headers and bad access tokens

	t.Run("test basic validation fails Without Header", func(t *testing.T) {
		_, err := client.GetUser("")
		if err == nil {
			t.Errorf("GetUser should have returned an error about the header")
		}
	})

	t.Run("test basic validation fails With Invalid Header", func(t *testing.T) {
		badAuthHeader := fmt.Sprintf("BadBearerHeader %s", accessToken)
		_, err := client.GetUser(badAuthHeader)
		if err == nil {
			t.Errorf("GetUser should have returned an error about the header")
		}
	})

	t.Run("test basic validation fails With Wrong Token", func(t *testing.T) {
		badAuthHeader := "Bearer thisisafaketoken"
		_, err := client.GetUser(badAuthHeader)
		if err == nil {
			t.Errorf("GetUser should have returned an error about the token")
		}
	})

	t.Run("test basic validation fails With Expired Token", func(t *testing.T) {
		// setup the expired token
		accessToken := testHelpers.CreateExpiredAccessToken(user, privateKey)
		authHeader := fmt.Sprintf("Bearer %s", accessToken)

		// run the test
		_, err := client.GetUser(authHeader)
		if err == nil {
			t.Errorf("GetUser should have returned an error about the token")
		}
	})

	t.Run("test basic validation fails With Bad Issuer", func(t *testing.T) {
		// setup the bad issuer
		tokenVerificationMetadataInput := &models.TokenVerificationMetadataInput{
			VerifierKey: publicKey,
			Issuer:      "newissuertestthatwontmatch",
		}

		client, err := propelauth.InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadataInput)
		if err != nil {
			t.Errorf("NewClient returned an error, cannot continue this test: %s", err)

			return
		}

		// run the test
		_, err = client.GetUser(authHeader)
		if err == nil {
			t.Errorf("GetUser should have returned an error about issuer")
		}
	})

	t.Run("test basic validation fails With Wrong Key", func(t *testing.T) {
		// generate a client with new random keys
		_, publicKey := testHelpers.GenerateRSAKeys()

		tokenVerificationMetadataInput := &models.TokenVerificationMetadataInput{
			VerifierKey: publicKey,
			Issuer:      "issuertest",
		}

		client, err := propelauth.InitBaseAuth("https://auth.example.com", "apikey", tokenVerificationMetadataInput)
		if err != nil {
			t.Errorf("NewClient returned an error, cannot even begin the tests: %s", err)

			return
		}

		// run the test
		_, err = client.GetUser(authHeader)
		if err == nil {
			t.Errorf("GetUser should have returned an error about decoding the token")
		}
	})
}
