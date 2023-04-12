package client

import (
	"fmt"
	"strconv"

	"encoding/json"
	"github.com/google/uuid"
	"net/url"
)

const urlPrefix = "https//propelauth.com/"
const backendUrlApiPrefix = "api/backend/v1/"

type Client struct {
	apiKey        string
	queryHelper   QueryHelper
	marshalHelper MarshalHelper
}

func NewClient(apiKey string) *Client {
	client := &Client{
		apiKey:        apiKey,
		queryHelper:   NewQueryHelper(urlPrefix, backendUrlApiPrefix),
		marshalHelper: NewMarshalHelper(),
	}

	return client
}

// public methods to fetch a user or users

func (o *Client) FetchUserMetadataByUserId(userId uuid.UUID, includeOrgs bool) (*UserMetadata, error) {
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

func (o *Client) FetchUserMetadataByEmail(email string, includeOrgs bool) (*UserMetadata, error) {
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

func (o *Client) FetchUserMetadataByUsername(username string, includeOrgs bool) (*UserMetadata, error) {
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

func (o *Client) FetchBatchUserMetadataByUserIds(userIds []string, includeOrgs bool) (*UserList, error) {
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

func (o *Client) FetchBatchUserMetadataByEmails(emails []string, includeOrgs bool) (*UserList, error) {
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

func (o *Client) FetchBatchUserMetadataByUsernames(usernames []string, includeOrgs bool) (*UserList, error) {
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

func (o *Client) FetchUsersByQuery(params UserQueryParams) (*UserList, error) {
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

func (o *Client) CreateUser(params CreateUserParams) (*UserID, error) {
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

func (o *Client) UpdateUserEmail(user_id uuid.UUID, params UpdateEmail) (bool, error) {
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

func (o *Client) UpdateUserMetadata(userId uuid.UUID, params UpdateUserMetadata) (bool, error) {
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

func (o *Client) UpdateUserPassword(userId uuid.UUID, params UpdateUserPasswordParam) (bool, error) {
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

func (o *Client) MigrateUserFromExternalSource(params MigrateUserParams) (bool, error) {
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

func (o *Client) FetchUsersInOrg(orgId uuid.UUID, params UserInOrgQueryParams) (*UserList, error) {
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

func (o *Client) AddUserToOrg(params AddUserToOrg) (bool, error) {
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

func (o *Client) FetchOrg(orgId uuid.UUID) (*OrgMetadata, error) {
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

func (o *Client) FetchOrgByQuery(params OrgQueryParams) (*OrgList, error) {
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

func (o *Client) CreateOrg(name string) (*OrgMetadata, error) {
	urlPostfix := "org"

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

func (o *Client) CreateAccessToken(userId uuid.UUID, durationInMinutes int) (*AccessToken, error) {
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

func (o *Client) CreateMagicLink(params CreateMagicLinkParams) (*CreateMagicLinkResponse, error) {
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

func (o *Client) ValidateAccessTokenAndGetUser(token string)                           {}
func (o *Client) ValidateAccessTokenAndGetUserserWithOrg(token string)                 {}
func (o *Client) ValidateAccessTokenAndGetUserserWithOrgByMinimumRole(token string)    {}
func (o *Client) ValidateAccessTokenAndGetUserserWithOrgByExactRole(token string)      {}
func (o *Client) ValidateAccessTokenAndGetUserserWithOrgByPermission(token string)     {}
func (o *Client) ValidateAccessTokenAndGetUserserWithOrgByAllPermissions(token string) {}

// private method to handle errors

func (o *Client) returnErrorMessageIfNotOk(queryResponse *QueryResponse) error {
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
