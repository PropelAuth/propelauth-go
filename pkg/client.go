package client

import (
	"fmt"
	"strconv"

	"encoding/json"
	"github.com/google/uuid"
	"net/url"

	"github.com/propelauth/propelauth-go/pkg/helpers"
	"github.com/propelauth/propelauth-go/pkg/models"
)

const urlPrefix = "https//propelauth.com/"
const backendUrlApiPrefix = "api/backend/v1/"

type Client struct {
	apiKey                    string
	authUrl                   string
	tokenVerificationMetadata models.TokenVerificationMetadata
	queryHelper               helpers.QueryHelperInterface
	marshalHelper             helpers.MarshalHelperInterface
	validationHelper          helpers.ValidationHelperInterface
}

func InitBaseAuth(authUrl string, apiKey string, tokenVerificationMetadata *models.TokenVerificationMetadata) (*Client, error) {
	// setup helpers
	queryHelper := helpers.NewQueryHelper(authUrl, backendUrlApiPrefix)
	marshalHelper := &helpers.MarshalHelper{}
	validationHelper := &helpers.ValidationHelper{}

	// validate the authUrl
	url, err := url.ParseRequestURI(authUrl)
	if err != nil {
		return nil, err
	} else if url.Scheme != "https" {
		return nil, fmt.Errorf("URL must start with https://")
	} else if url.Path != "" {
		return nil, fmt.Errorf("URL must not end with a trailing slash")
	} else if url.Host == "" {
		return nil, fmt.Errorf("Invalid URL")
	}

	// if tokenVerificationMetadata wasn't passed in, create one
	if tokenVerificationMetadata == nil {
		endpointURL := "https://" + url.Host + "/api/v1/token_verification_metadata"

		queryResponse, err := queryHelper.RequestHelper("GET", apiKey, endpointURL, nil)
		if err != nil {
			return nil, err
		}

		if queryResponse.StatusCode == 401 {
			return nil, fmt.Errorf("apiKey is incorrect")
		} else if queryResponse.StatusCode == 400 {
			return nil, fmt.Errorf("Bad request: %s", queryResponse.ResponseText)
		} else if queryResponse.StatusCode == 404 {
			return nil, fmt.Errorf("URL is incorrect")
		} else if queryResponse.StatusCode != 200 { // this must be last
			return nil, fmt.Errorf("Unknown error when fetching token verification metadata")
		}

		authTokenVerificationMetadataResponse, err := marshalHelper.GetAuthTokenVerificationMetadataResponseFromBytes(queryResponse.BodyBytes)
		if err != nil {
			return nil, err
		}

		tokenVerificationMetadata = &models.TokenVerificationMetadata{
			VerifierKey: authTokenVerificationMetadataResponse.PublicKeyPem,
			Issuer:      authUrl,
		}
	}

	client := &Client{
		apiKey:                    apiKey,
		authUrl:                   authUrl,
		tokenVerificationMetadata: *tokenVerificationMetadata,
		queryHelper:               queryHelper,
		marshalHelper:             marshalHelper,
		validationHelper:          validationHelper,
	}

	return client, nil
}

// public methods to fetch a user or users

func (o *Client) FetchUserMetadataByUserId(userId uuid.UUID, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := fmt.Sprintf("user/%s", userId)

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.apiKey, urlPostfix, queryParams)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	user, err := o.marshalHelper.GetUserMetadataFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (o *Client) FetchUserMetadataByEmail(email string, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := "user/email"

	queryParams := url.Values{
		"email":        {email},
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.apiKey, urlPostfix, queryParams)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	user, err := o.marshalHelper.GetUserMetadataFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (o *Client) FetchUserMetadataByUsername(username string, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := "user/username"

	queryParams := url.Values{
		"username":     {username},
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.apiKey, urlPostfix, queryParams)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	user, err := o.marshalHelper.GetUserMetadataFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (o *Client) FetchBatchUserMetadataByUserIds(userIds []string, includeOrgs bool) (*models.UserList, error) {
	urlPostfix := "user/user_ids"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	type UserIds struct {
		UserIds []string `json:"user_ids"`
	}

	bodyParams := UserIds{
		UserIds: userIds,
	}

	bodyJson, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, queryParams, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	users, err := o.marshalHelper.GetUserListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (o *Client) FetchBatchUserMetadataByEmails(emails []string, includeOrgs bool) (*models.UserList, error) {
	urlPostfix := "user/emails"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}
	type UserIds struct {
		Emails []string `json:"emails"`
	}

	bodyParams := UserIds{
		Emails: emails,
	}

	bodyJson, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, queryParams, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	users, err := o.marshalHelper.GetUserListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (o *Client) FetchBatchUserMetadataByUsernames(usernames []string, includeOrgs bool) (*models.UserList, error) {
	urlPostfix := "user/usernames"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	type UserIds struct {
		Usernames []string `json:"usernames"`
	}

	bodyParams := UserIds{
		Usernames: usernames,
	}

	bodyJson, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, queryParams, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	users, err := o.marshalHelper.GetUserListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (o *Client) FetchUsersByQuery(params models.UserQueryParams) (*models.UserList, error) {
	urlPostfix := "user/query"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	users, err := o.marshalHelper.GetUserListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// public methods to modify a user

func (o *Client) CreateUser(params models.CreateUserParams) (*models.UserID, error) {
	urlPostfix := "user"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	user, err := o.marshalHelper.GetUserFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (o *Client) UpdateUserEmail(user_id uuid.UUID, params models.UpdateEmail) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/email", user_id)

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) UpdateUserMetadata(userId uuid.UUID, params models.UpdateUserMetadata) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/metadata", userId)

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) UpdateUserPassword(userId uuid.UUID, params models.UpdateUserPasswordParam) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/password", userId)

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) MigrateUserFromExternalSource(params models.MigrateUserParams) (bool, error) {
	urlPostfix := "migrate_user"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) DeleteUser(userId uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s", userId)

	queryResponse, err := o.queryHelper.Delete(o.apiKey, urlPostfix, nil)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) DisableUser(userId uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/disable", userId)

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, nil)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) EnableUser(userId uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/enable", userId)

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, nil)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

// public methods for users in orgs

func (o *Client) FetchUsersInOrg(orgId uuid.UUID, params models.UserInOrgQueryParams) (*models.UserList, error) {
	urlPostfix := fmt.Sprintf("user/org/%s", orgId)

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	users, err := o.marshalHelper.GetUserListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (o *Client) AddUserToOrg(params models.AddUserToOrg) (bool, error) {
	urlPostfix := "org/add_user"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

// public methods for orgs

func (o *Client) FetchOrg(orgId uuid.UUID) (*models.OrgMetadata, error) {
	urlPostfix := fmt.Sprintf("org/%s", orgId)

	queryResponse, err := o.queryHelper.Get(o.apiKey, urlPostfix, nil)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	org, err := o.marshalHelper.GetOrgMetadataFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (o *Client) FetchOrgByQuery(params models.OrgQueryParams) (*models.OrgList, error) {
	urlPostfix := "org/query"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	orgs, err := o.marshalHelper.GetOrgListFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (o *Client) CreateOrg(name string) (*models.OrgMetadata, error) {
	urlPostfix := "org"

	// assemble the parameters

	type CreateOrg struct {
		Name string `json:"name"`
	}

	bodyParams := CreateOrg{
		Name: name,
	}

	bodyJson, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	org, err := o.marshalHelper.GetOrgMetadataFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (o *Client) AllowOrgToSetupSamlConnection(orgId uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s/allow_saml", orgId)

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, nil)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

func (o *Client) DisallowOrgToSetupSamlConnection(orgId uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s/disallow_saml", orgId)

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, nil)
	if err != nil {
		return false, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, err
	}

	return true, nil
}

// public methods for misc functionality

func (o *Client) CreateAccessToken(userId uuid.UUID, durationInMinutes int) (*models.AccessToken, error) {
	urlPostfix := "access_token"

	// assemble body params

	type CreateAccessToken struct {
		UserID            uuid.UUID `json:"user_id"`
		DurationInMinutes int       `json:"durationInMinutes"`
	}

	bodyParams := CreateAccessToken{
		UserID:            userId,
		DurationInMinutes: durationInMinutes,
	}

	bodyJson, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if queryResponse.StatusCode == 401 {
		return nil, fmt.Errorf("API Key is incorrect")
	} else if queryResponse.StatusCode == 400 {
		return nil, fmt.Errorf("Bad request: %s", queryResponse.ResponseText)
	} else if queryResponse.StatusCode == 403 {
		return nil, fmt.Errorf("User not found")
	} else if queryResponse.StatusCode == 404 {
		return nil, fmt.Errorf("Access token creation not enabled")
	} else if queryResponse.StatusCode != 200 { // this must be last
		return nil, fmt.Errorf("Unknown error when performing operation")
	}

	accessToken, err := o.marshalHelper.GetAccessTokenFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (o *Client) CreateMagicLink(params models.CreateMagicLinkParams) (*models.CreateMagicLinkResponse, error) {
	urlPostfix := "magic_link"

	bodyJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	queryResponse, err := o.queryHelper.Post(o.apiKey, urlPostfix, nil, bodyJson)
	if err != nil {
		return nil, err
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, err
	}

	magicLink, err := o.marshalHelper.GetMagicLinkResponseFromBytes(queryResponse.BodyBytes)
	if err != nil {
		return nil, err
	}

	return magicLink, nil
}

// public methods around authorization

func (o *Client) ValidateAccessTokenAndGetUser(authHeader string) (*models.UserFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (o *Client) ValidateAccessTokenAndGetUserWithOrg(authHeader string, orgId uuid.UUID) (*models.UserAndOrgMemberInfoFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	orgMemberInfo, err := o.validationHelper.ValidateOrgAccessAndGetOrgMemberInfo(user, orgId)
	if err != nil {
		return nil, err
	}

	return &models.UserAndOrgMemberInfoFromToken{User: *user, OrgMemberInfo: *orgMemberInfo}, nil
}

func (o *Client) ValidateAccessTokenAndGetUserWithOrgByMinimumRole(authHeader string, orgId uuid.UUID, minimumRole string) (*models.UserAndOrgMemberInfoFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	orgMemberInfo, err := o.validationHelper.ValidateOrgAccessAndGetOrgMemberInfoByMinimumRole(user, orgId, minimumRole)
	if err != nil {
		return nil, err
	}

	return &models.UserAndOrgMemberInfoFromToken{User: *user, OrgMemberInfo: *orgMemberInfo}, nil
}

func (o *Client) ValidateAccessTokenAndGetUserWithOrgByExactRole(authHeader string, orgId uuid.UUID, exactRole string) (*models.UserAndOrgMemberInfoFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	orgMemberInfo, err := o.validationHelper.ValidateOrgAccessAndGetOrgMemberInfoByExactRole(user, orgId, exactRole)
	if err != nil {
		return nil, err
	}

	return &models.UserAndOrgMemberInfoFromToken{User: *user, OrgMemberInfo: *orgMemberInfo}, nil
}

func (o *Client) ValidateAccessTokenAndGetUserWithOrgByPermission(authHeader string, orgId uuid.UUID, permission string) (*models.UserAndOrgMemberInfoFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	orgMemberInfo, err := o.validationHelper.ValidateOrgAccessAndGetOrgMemberInfoByPermission(user, orgId, permission)
	if err != nil {
		return nil, err
	}

	return &models.UserAndOrgMemberInfoFromToken{User: *user, OrgMemberInfo: *orgMemberInfo}, nil
}

func (o *Client) ValidateAccessTokenAndGetUserWithOrgByAllPermissions(authHeader string, orgId uuid.UUID, permissions []string) (*models.UserAndOrgMemberInfoFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, err
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, err
	}

	orgMemberInfo, err := o.validationHelper.ValidateOrgAccessAndGetOrgMemberInfoByAllPermissions(user, orgId, permissions)
	if err != nil {
		return nil, err
	}

	return &models.UserAndOrgMemberInfoFromToken{User: *user, OrgMemberInfo: *orgMemberInfo}, nil
}

// private method to handle errors

func (o *Client) returnErrorMessageIfNotOk(queryResponse *helpers.QueryResponse) error {
	if queryResponse.StatusCode == 401 {
		return fmt.Errorf("API Key is incorrect")
	} else if queryResponse.StatusCode == 400 {
		return fmt.Errorf("Bad request: %s", queryResponse.ResponseText)
	} else if queryResponse.StatusCode == 404 {
		return fmt.Errorf("API not found")
	} else if queryResponse.StatusCode == 426 {
		return fmt.Errorf("Cannot use organizations unless B2B support is enabled. Enable it in your PropelAuth dashboard.")
	} else if queryResponse.StatusCode == 429 {
		return fmt.Errorf("Your app is making too many requests, too quickly")
	} else if queryResponse.StatusCode != 200 { // this must be last
		return fmt.Errorf("Unknown error when performing operation")
	}

	return nil
}
