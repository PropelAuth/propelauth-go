// Package client is the main package for the PropelAuth Go library.
package client

import (
	"fmt"
	"strconv"

	"encoding/json"
	"net/url"

	"github.com/google/uuid"

	"github.com/propelauth/propelauth-go/pkg/helpers"
	"github.com/propelauth/propelauth-go/pkg/models"
)

const backendURLApiPrefix = "/api/backend/v1/"

// ClientInterface contains all the methods for interacting with the PropelAuth backend and the JWT.
// It's also a convient listing of all the methods available to the integration programmer.
type ClientInterface interface {
	// user endpoints
	CreateAccessToken(userID uuid.UUID, durationInMinutes int, createAccessTokenOptions ...models.CreateAccessTokenOptions) (*models.AccessToken, error)
	CreateMagicLink(params models.CreateMagicLinkParams) (*models.CreateMagicLinkResponse, error)
	CreateUser(params models.CreateUserParams) (*models.UserID, error)
	DeleteUser(userID uuid.UUID) (bool, error)
	DisableUser(userID uuid.UUID) (bool, error)
	EnableUser(userID uuid.UUID) (bool, error)
	FetchBatchUserMetadataByEmails(emails []string, includeOrgs bool) (map[string]models.UserMetadata, error)
	FetchBatchUserMetadataByUserIds(userIds []uuid.UUID, includeOrgs bool) (map[uuid.UUID]models.UserMetadata, error)
	FetchBatchUserMetadataByUsernames(usernames []string, includeOrgs bool) (map[string]models.UserMetadata, error)
	FetchUserMetadataByEmail(email string, includeOrgs bool) (*models.UserMetadata, error)
	FetchUserMetadataByUserID(userID uuid.UUID, includeOrgs bool) (*models.UserMetadata, error)
	FetchUserMetadataByUsername(username string, includeOrgs bool) (*models.UserMetadata, error)
	FetchUsersByQuery(params models.UserQueryParams) (*models.UserList, error)
	MigrateUserFromExternalSource(params models.MigrateUserParams) (bool, error)
	UpdateUserEmail(userID uuid.UUID, params models.UpdateEmail) (bool, error)
	UpdateUserMetadata(userID uuid.UUID, params models.UpdateUserMetadata) (bool, error)
	UpdateUserPassword(userID uuid.UUID, params models.UpdateUserPasswordParam) (bool, error)
	EnableUserCanCreateOrgs(userID uuid.UUID) (bool, error)
	DisableUserCanCreateOrgs(userID uuid.UUID) (bool, error)
	ClearUserPassword(userID uuid.UUID) (bool, error)
	DisableUser2fa(userID uuid.UUID) (bool, error)
	ResendEmailConfirmation(userID uuid.UUID) (bool, error)
	LogoutAllUserSessions(userID uuid.UUID) (bool, error)

	// org endpoints
	AllowOrgToSetupSamlConnection(orgID uuid.UUID) (bool, error)
	CreateOrg(name string) (*models.OrgMetadata, error)
	CreateOrgV2(params models.CreateOrgV2Params) (*models.CreateOrgV2Response, error)
	DeleteOrg(orgID uuid.UUID) (bool, error)
	DisallowOrgToSetupSamlConnection(orgID uuid.UUID) (bool, error)
	FetchOrg(orgID uuid.UUID) (*models.OrgMetadata, error)
	FetchOrgByQuery(params models.OrgQueryParams) (*models.OrgList, error)
	FetchCustomRoleMappings() (*models.CustomRoleMappingList, error)
	FetchPendingInvites(params models.FetchPendingInvitesParams) (*models.PendingInvitesPage, error)
	UpdateOrgMetadata(orgID uuid.UUID, params models.UpdateOrg) (bool, error)
	SubscribeOrgToRoleMapping(orgID uuid.UUID, params models.OrgRoleMappingSubscription) (bool, error)
	ChangeUserRoleInOrg(params models.ChangeUserRoleInOrg) (bool, error)
	RevokePendingOrgInvite(params models.RevokePendingOrgInvite) (bool, error)
	CreateOrgSamlConnectionLink(orgID uuid.UUID, params models.CreateSamlConnectionLinkBody) (*models.CreateSamlConnectionLinkResponse, error)
	FetchSamlSpMetadata(orgID uuid.UUID) (*models.SamlSpMetadata, error)
	SetSamlIdpMetadata(params models.SamlIdpMetadata) (bool, error)
	SamlGoLive(orgId uuid.UUID) (bool, error)
	DeleteSamlConnection(orgId uuid.UUID) (bool, error)

	// user in org endpoints
	AddUserToOrg(params models.AddUserToOrg) (bool, error)
	RemoveUserFromOrg(params models.RemoveUserFromOrg) (bool, error)
	InviteUserToOrg(params models.InviteUserToOrg) (bool, error)
	FetchUsersInOrg(orgID uuid.UUID, params models.UserInOrgQueryParams) (*models.UserList, error)

	// api key endpoints
	FetchAPIKey(apiKeyID string) (*models.APIKeyFull, error)
	CreateAPIKey(params models.APIKeyCreateParams) (*models.APIKeyNew, error)
	UpdateAPIKey(apiKeyID string, params models.APIKeyUpdateParams) (bool, error)
	DeleteAPIKey(apiKeyID string) (bool, error)
	FetchCurrentAPIKeys(params models.APIKeysQueryParams) (*models.APIKeyResultPage, error)
	FetchArchivedAPIKeys(params models.APIKeysQueryParams) (*models.APIKeyResultPage, error)
	ValidatePersonalAPIKey(apiKeyToken string) (*models.PersonalAPIKeyValidation, error)
	ValidateOrgAPIKey(apiKeyToken string) (*models.OrgAPIKeyValidation, error)
	ValidateAPIKey(apiKeyToken string) (*models.APIKeyValidation, error)

	// a method to validate the JWT
	GetUser(authHeader string) (*models.UserFromToken, error)
}

// Client is the main struct for the PropelAuth Go library. It contains all the methods for interacting with the
// PropelAuth backend.
type Client struct {
	integrationAPIKey         string
	authURL                   string
	tokenVerificationMetadata models.TokenVerificationMetadata
	queryHelper               helpers.QueryHelperInterface
	validationHelper          helpers.ValidationHelperInterface
}

// InitBaseAuth initializes the PropelAuth client with the authURL, integrationAPIKey.
//
// This is the normal entrance to accessing the PropelAuth backend.
//
// The authURL and integrationAPIKey can be found in your PropelAuth dashboard, in the "Backend Integrations" section.
// You can pass in a tokenVerificationMetadata if you have it, but it's not required.
func InitBaseAuth(authURL string, integrationAPIKey string, tokenVerificationMetadataInput *models.TokenVerificationMetadataInput) (ClientInterface, error) {
	// setup helpers
	queryHelper := helpers.NewQueryHelper(authURL, backendURLApiPrefix)
	validationHelper := &helpers.ValidationHelper{}

	// validate the authURL
	url, err := url.ParseRequestURI(authURL)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse the authURL: %w", err)
	} else if url.Scheme != "https" {
		return nil, fmt.Errorf("URL must start with https://")
	} else if url.Path != "" {
		return nil, fmt.Errorf("URL must not end with a trailing slash")
	} else if url.Host == "" {
		return nil, fmt.Errorf("Invalid URL")
	}

	var tokenVerificationMetadata *models.TokenVerificationMetadata

	// if tokenVerificationMetadata wasn't passed in, create one
	if tokenVerificationMetadataInput == nil {
		endpointURL := "https://" + url.Host + "/api/v1/token_verification_metadata"

		queryResponse, err := queryHelper.RequestHelper("GET", integrationAPIKey, endpointURL, nil)
		if err != nil {
			return nil, fmt.Errorf("Error on fetching token verification metadata: %w", err)
		}

		if queryResponse.StatusCode != 200 {
			switch statusCode := queryResponse.StatusCode; statusCode {
			case 401:
				return nil, fmt.Errorf("integrationAPIKey is incorrect")
			case 400:
				return nil, fmt.Errorf("Bad request: %s", queryResponse.BodyText)
			case 404:
				return nil, fmt.Errorf("URL is incorrect")
			default:
				return nil, fmt.Errorf("Unknown error when fetching token verification metadata. Status code: %s. Body: %s", strconv.Itoa(queryResponse.StatusCode), queryResponse.BodyText)
			}
		}

		authTokenVerificationMetadataResponse := &models.AuthTokenVerificationMetadataResponse{}
		if err := json.Unmarshal(queryResponse.BodyBytes, authTokenVerificationMetadataResponse); err != nil {
			return nil, fmt.Errorf("Error on unmarshalling bytes to AuthTokenVerificationMetadataResponse: %w", err)
		}

		rsaPublicKey, err := validationHelper.ConvertPEMStringToRSAPublicKey(authTokenVerificationMetadataResponse.VerifierKeyPem)
		if err != nil {
			return nil, fmt.Errorf("Error converting a PEM string to an RSA Public Key: %w", err)
		}

		tokenVerificationMetadata = &models.TokenVerificationMetadata{
			VerifierKey: *rsaPublicKey,
			Issuer:      authURL,
		}
	} else {
		rsaPublicKey, err := validationHelper.ConvertPEMStringToRSAPublicKey(tokenVerificationMetadataInput.VerifierKey)
		if err != nil {
			return nil, fmt.Errorf("Error converting a PEM string to an RSA Public Key: %w", err)
		}

		tokenVerificationMetadata = &models.TokenVerificationMetadata{
			VerifierKey: *rsaPublicKey,
			Issuer:      tokenVerificationMetadataInput.Issuer,
		}
	}

	client := &Client{
		integrationAPIKey:         integrationAPIKey,
		authURL:                   authURL,
		tokenVerificationMetadata: *tokenVerificationMetadata,
		queryHelper:               queryHelper,
		validationHelper:          validationHelper,
	}

	return client, nil
}

// Public methods to fetch a user or users

// FetchUserMetadataByUserID will fetch a single user by their user ID. If includeOrgs is true, we'll also
// fetch the organizations data for each organization the user is in.
func (o *Client) FetchUserMetadataByUserID(userID uuid.UUID, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := fmt.Sprintf("user/%s", userID)

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching user by id: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching user by id: %w", err)
	}

	user := &models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserMetadata: %w", err)
	}

	return user, nil
}

// FetchUserMetadataByEmail will fetch a single user by their email. If includeOrgs is true, we'll also
// fetch the organizations data for each organization the user is in.
func (o *Client) FetchUserMetadataByEmail(email string, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := "user/email"

	queryParams := url.Values{
		"email":        {email},
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching user by email: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching user by email: %w", err)
	}

	user := &models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserMetadata: %w", err)
	}

	return user, nil
}

// FetchUserMetadataByUsername will fetch a single user by their username. If includeOrgs is true, we'll also
// fetch the organizations data for each organization the user is in.
func (o *Client) FetchUserMetadataByUsername(username string, includeOrgs bool) (*models.UserMetadata, error) {
	urlPostfix := "user/username"

	queryParams := url.Values{
		"username":     {username},
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching user by username: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching user by username: %w", err)
	}

	user := &models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserMetadata: %w", err)
	}

	return user, nil
}

// FetchBatchUserMetadataByUserIds will fetch all the users with the listed IDS. If includeOrgs is true, we'll
// also fetch the organizations data for each organization the user is in.
func (o *Client) FetchBatchUserMetadataByUserIds(userIds []uuid.UUID, includeOrgs bool) (map[uuid.UUID]models.UserMetadata, error) {
	urlPostfix := "user/user_ids"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	type UserIds struct {
		UserIds []uuid.UUID `json:"user_ids"`
	}

	bodyParams := UserIds{
		UserIds: userIds,
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, queryParams, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by ids: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by ids: %w", err)
	}

	users := &[]models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to []UserMetadata: %w", err)
	}

	// assemble the return

	userMap := map[uuid.UUID]models.UserMetadata{}
	for _, user := range *users {
		userMap[user.UserID] = user
	}

	return userMap, nil
}

// FetchBatchUserMetadataByEmails will fetch all the users with the listed emails. If includeOrgs is true, we'll
// also fetch the organizations data for each organization the user is in.
func (o *Client) FetchBatchUserMetadataByEmails(emails []string, includeOrgs bool) (map[string]models.UserMetadata, error) {
	urlPostfix := "user/emails"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	type UserEmails struct {
		Emails []string `json:"emails"`
	}

	bodyParams := UserEmails{
		Emails: emails,
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, queryParams, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by emails: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by emails: %w", err)
	}

	users := &[]models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to []UserMetadata: %w", err)
	}

	// assemble the return

	userMap := map[string]models.UserMetadata{}
	for _, user := range *users {
		userMap[user.Email] = user
	}

	return userMap, nil
}

// FetchBatchUserMetadataByUsernames will fetch all the users with the listed usernames. If includeOrgs is true,
// we'll also fetch the organizations data for each organization the user is in.
func (o *Client) FetchBatchUserMetadataByUsernames(usernames []string, includeOrgs bool) (map[string]models.UserMetadata, error) {
	urlPostfix := "user/usernames"

	// assemble the parameters

	queryParams := url.Values{
		"include_orgs": {strconv.FormatBool(includeOrgs)},
	}

	type UserNames struct {
		Usernames []string `json:"usernames"`
	}

	bodyParams := UserNames{
		Usernames: usernames,
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, queryParams, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by usernames: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching batch users by usernames: %w", err)
	}

	users := &[]models.UserMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to []UserMetadata: %w", err)
	}

	// assemble the return

	userMap := map[string]models.UserMetadata{}
	for _, user := range *users {
		userMap[*user.Username] = user
	}

	return userMap, nil
}

// FetchUsersByQuery will fetch a paged list of users.
func (o *Client) FetchUsersByQuery(params models.UserQueryParams) (*models.UserList, error) {
	urlPostfix := "user/query"

	queryParams := url.Values{}

	if params.PageNumber != nil {
		queryParams.Add("page_number", strconv.Itoa(*params.PageNumber))
	}
	if params.PageSize != nil {
		queryParams.Add("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.OrderBy != nil {
		queryParams.Add("order_by", *params.OrderBy)
	}
	if params.EmailOrUsername != nil {
		queryParams.Add("email_or_username", *params.EmailOrUsername)
	}
	if params.IncludeOrgs != nil {
		queryParams.Add("include_orgs", strconv.FormatBool(*params.IncludeOrgs))
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching users by query: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching users by query: %w", err)
	}

	users := &models.UserList{}
	if err := json.Unmarshal(queryResponse.BodyBytes, users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserList: %w", err)
	}

	return users, nil
}

// public methods to modify a user

// CreateUser will create a new user.
func (o *Client) CreateUser(params models.CreateUserParams) (*models.UserID, error) {
	urlPostfix := "user/"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating user: %w", err)
	}

	user := &models.UserID{}
	if err := json.Unmarshal(queryResponse.BodyBytes, user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserID: %w", err)
	}

	return user, nil
}

// UpdateUserEmail will update a user's email address. if RequireEmailConfirmation is set to true, we'll send
// out an email to confirm the new email address.
func (o *Client) UpdateUserEmail(userID uuid.UUID, params models.UpdateEmail) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/email", userID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on updating user email: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on updating user email: %w", err)
	}

	return true, nil
}

// ClearUserPassword will clear a user's password.
func (o *Client) ClearUserPassword(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/clear_password", userID)

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on clearing user password: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on clearing user password: %w", err)
	}

	return true, nil
}

// UpdateUserMetadata will update properties on a user. All fields are optional, we'll only update the ones
// that are provided.
func (o *Client) UpdateUserMetadata(userID uuid.UUID, params models.UpdateUserMetadata) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s", userID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on updating user metadata: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on updating user metadata: %w", err)
	}

	return true, nil
}

// UpdateUserPassword will update a user's password. If the user is logged in, they will be logged out.
func (o *Client) UpdateUserPassword(userID uuid.UUID, params models.UpdateUserPasswordParam) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/password", userID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on updating user password: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on updating user password: %w", err)
	}

	return true, nil
}

// MigrateUserFromExternalSource will migrate a user from another system.
func (o *Client) MigrateUserFromExternalSource(params models.MigrateUserParams) (bool, error) {
	urlPostfix := "migrate_user/"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on migrating user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on migrating user: %w", err)
	}

	return true, nil
}

// DeleteUser will delete a user, removing them from all organizations they are in. This is a permanent
// action and cannot be undone. If you're unsure if you want this, use DisableUser instead.
func (o *Client) DeleteUser(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s", userID)

	queryResponse, err := o.queryHelper.Delete(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on deleting user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on deleting user: %w", err)
	}

	return true, nil
}

// DisableUser will disable a user, preventing them from logging in.
func (o *Client) DisableUser(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/disable", userID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on disabling user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on disabling user: %w", err)
	}

	return true, nil
}

// EnableUser will enable a user, meaning they will be allowed to logging in.
func (o *Client) EnableUser(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/enable", userID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on enabling user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on enabling user: %w", err)
	}

	return true, nil
}

// EnableUserCanCreateOrgs will let a user create orgs even when the global users_can_create_orgs is set to false.
func (o *Client) EnableUserCanCreateOrgs(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/can_create_orgs/enable", userID)

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on enable user can create orgs: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on enable user can create orgs: %w", err)
	}

	return true, nil
}

// DisableUserCanCreateOrgs will prevent a user from creating orgs, unless the global users_can_create_orgs is set to true.
func (o *Client) DisableUserCanCreateOrgs(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/can_create_orgs/disable", userID)

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on disable user can create orgs: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on disable user can create orgs: %w", err)
	}

	return true, nil
}

// DisableUser2fa will disable 2fa for a user.
func (o *Client) DisableUser2fa(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/disable_2fa", userID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on disabling user 2fa: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on disabling user 2fa: %w", err)
	}

	return true, nil
}

// public methods for users in orgs

// FetchUsersInOrg will fetch a paged list of users in an organization.
func (o *Client) FetchUsersInOrg(orgID uuid.UUID, params models.UserInOrgQueryParams) (*models.UserList, error) {
	urlPostfix := fmt.Sprintf("user/org/%s", orgID)

	queryParams := url.Values{}

	if params.PageNumber != nil {
		queryParams.Add("page_number", strconv.Itoa(*params.PageNumber))
	}
	if params.PageSize != nil {
		queryParams.Add("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.IncludeOrgs != nil {
		queryParams.Add("include_orgs", strconv.FormatBool(*params.IncludeOrgs))
	}
	if params.Role != nil {
		queryParams.Add("role", *params.Role)
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching users in org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching users in org: %w", err)
	}

	users := &models.UserList{}
	if err := json.Unmarshal(queryResponse.BodyBytes, users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserList: %w", err)
	}

	return users, nil
}

// AddUserToOrg will add a user to an org with a role.
func (o *Client) AddUserToOrg(params models.AddUserToOrg) (bool, error) {
	urlPostfix := "org/add_user"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on adding user to org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on adding user to org: %w", err)
	}

	return true, nil
}

// RemoveUserFromOrg will remove the user from an org.
func (o *Client) RemoveUserFromOrg(params models.RemoveUserFromOrg) (bool, error) {
	urlPostfix := "org/remove_user"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on removing user from org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on removing user from org: %w", err)
	}

	return true, nil
}

// ChangeUserRole will change a user's role in an org.
func (o *Client) ChangeUserRoleInOrg(params models.ChangeUserRoleInOrg) (bool, error) {
	urlPostfix := "org/change_role"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on changing user role in org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on changing user role in org: %w", err)
	}

	return true, nil
}

// InviteUserToOrg will email a user and invite them to join an org. If they don't have an account
//
//	yet, they'll be asked to make one, and will be able to join the org right afterwards.
func (o *Client) InviteUserToOrg(params models.InviteUserToOrg) (bool, error) {
	urlPostfix := "invite_user"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on inviting user to org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on inviting user to org: %w", err)
	}

	return true, nil
}

// ResendEmailConfirmation will resend the email confirmation email to a user.
func (o *Client) ResendEmailConfirmation(userID uuid.UUID) (bool, error) {
	urlPostfix := "resend_email_confirmation"

	type ResendEmailConfirmationParams struct {
		UserID uuid.UUID `json:"user_id"`
	}

	bodyParams := ResendEmailConfirmationParams{
		UserID: userID,
	}
	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on resending email confirmation to user: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on resending email confirmation to user: %w", err)
	}

	return true, nil
}

// LogoutAllUserSessions will log out all of a user's sessions.
func (o *Client) LogoutAllUserSessions(userID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("user/%s/logout_all_sessions", userID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on logging out all user sessions : %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on logging out all user sessions: %w", err)
	}

	return true, nil
}

// public methods for orgs

// FetchOrg will fetch an org's data.
func (o *Client) FetchOrg(orgID uuid.UUID) (*models.OrgMetadata, error) {
	urlPostfix := fmt.Sprintf("org/%s", orgID)

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, nil)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching org: %w", err)
	}

	org := &models.OrgMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, org); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to OrgMetadata: %w", err)
	}

	return org, nil
}

// FetchOrgByQuery will fetch a paged list of organizations.
func (o *Client) FetchOrgByQuery(params models.OrgQueryParams) (*models.OrgList, error) {
	urlPostfix := "org/query"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching orgs by query: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching orgs by query: %w", err)
	}

	orgs := &models.OrgList{}
	if err := json.Unmarshal(queryResponse.BodyBytes, orgs); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to OrgList: %w", err)
	}

	return orgs, nil
}

// FetchCustomRoleMappings will fetch all custom Role-to-Permissions mappings available.
func (o *Client) FetchCustomRoleMappings() (*models.CustomRoleMappingList, error) {
	urlPostfix := "custom_role_mappings"

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, nil)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching custom_role_mappings: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching custom_role_mappings: %w", err)
	}

	custom_role_mappings := &models.CustomRoleMappingList{}
	if err := json.Unmarshal(queryResponse.BodyBytes, custom_role_mappings); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to CustomRoleMappingList: %w", err)
	}

	return custom_role_mappings, nil
}

// FetchPendingInvites will fetch a paged list of pending invites.
func (o *Client) FetchPendingInvites(params models.FetchPendingInvitesParams) (*models.PendingInvitesPage, error) {
	urlPostfix := "pending_org_invites"

	queryParams := url.Values{}

	if params.PageNumber != nil {
		queryParams.Add("page_number", strconv.Itoa(*params.PageNumber))
	}
	if params.PageSize != nil {
		queryParams.Add("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.OrgID != nil {
		queryParams.Add("org_id", params.OrgID.String())
	}

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching pending invites: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching pending invites: %w", err)
	}

	pendingInvites := &models.PendingInvitesPage{}
	if err := json.Unmarshal(queryResponse.BodyBytes, pendingInvites); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to PendingInvitesPage: %w", err)
	}

	return pendingInvites, nil
}

func (o *Client) RevokePendingOrgInvite(params models.RevokePendingOrgInvite) (bool, error) {
	urlPostfix := "pending_org_invites"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Delete(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error deleting pending org invite: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error deleting pending org invite: %w", err)
	}

	return true, nil
}

// NOTE: THIS IS DEPRECATED.
// CreateOrg will an organization and return its data, which is mostly just the org's ID.
func (o *Client) CreateOrg(name string) (*models.OrgMetadata, error) {
	urlPostfix := "org/"

	// assemble the parameters

	type CreateOrg struct {
		Name string `json:"name"`
	}

	bodyParams := CreateOrg{
		Name: name,
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating org: %w", err)
	}

	org := &models.OrgMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, org); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to OrgMetadata: %w", err)
	}

	return org, nil
}

// CreateOrgV2 is the updated version of CreateOrg. It will create an organization and return its data.
func (o *Client) CreateOrgV2(params models.CreateOrgV2Params) (*models.CreateOrgV2Response, error) {
	urlPostfix := "org/"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating org: %w", err)
	}

	org := &models.CreateOrgV2Response{}
	if err := json.Unmarshal(queryResponse.BodyBytes, org); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to CreateOrgV2Response: %w", err)
	}

	return org, nil
}

// DeleteOrg will delete an organization.
func (o *Client) DeleteOrg(orgID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s", orgID)

	queryResponse, err := o.queryHelper.Delete(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on deleting an org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on deleting an org: %w", err)
	}

	return true, nil
}

// UpdateOrgMetadata will update properties on an organization. All fields are optional, we'll only update the ones
// that are provided.
func (o *Client) UpdateOrgMetadata(orgID uuid.UUID, params models.UpdateOrg) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s", orgID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on updating org metadata: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on updating org metadata: %w", err)
	}

	return true, nil
}

// SubscribeOrgToRoleMapping will subscribe the organization to a role mapping.
func (o *Client) SubscribeOrgToRoleMapping(orgID uuid.UUID, params models.OrgRoleMappingSubscription) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s", orgID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Put(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on subscribing org to a role mapping: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on subscribing org to a role mapping: %w", err)
	}

	return true, nil
}

// AllowOrgToSetupSamlConnection will turn on an org's ability to setup a SAML connection.
func (o *Client) AllowOrgToSetupSamlConnection(orgID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s/allow_saml", orgID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on allowing org to setup SAML connection: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on allowing org to setup SAML connection: %w", err)
	}

	return true, nil
}

// DisallowOrgToSetupSamlConnection will turn off an org's ability to setup a SAML connection. This is the default.
func (o *Client) DisallowOrgToSetupSamlConnection(orgID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("org/%s/disallow_saml", orgID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on disallowing org to setup SAML connection: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on disallowing org to setup SAML connection: %w", err)
	}

	return true, nil
}

// CreateOrgSamlConnectionLink will create a SAML connection link for an org.
func (o *Client) CreateOrgSamlConnectionLink(orgID uuid.UUID, params models.CreateSamlConnectionLinkBody) (*models.CreateSamlConnectionLinkResponse, error) {
	urlPostfix := fmt.Sprintf("org/%s/create_saml_connection_link", orgID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating SAML connection link for org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating SAML connection link for org: %w", err)
	}

	newUrl := &models.CreateSamlConnectionLinkResponse{}
	if err := json.Unmarshal(queryResponse.BodyBytes, newUrl); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to CreateSamlConnectionLinkResponse: %w", err)
	}

	return newUrl, nil
}

// FetchSamlSpMetadata will fetch the Service Provider metadata needed to configure a SAML connection for an org.
func (o *Client) FetchSamlSpMetadata(orgID uuid.UUID) (*models.SamlSpMetadata, error) {
	urlPostfix := fmt.Sprintf("saml_sp_metadata/%s", orgID)

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, nil)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching SAML SP Metadata for org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching SAML SP Metadata for org: %w", err)
	}

	spMetadata := &models.SamlSpMetadata{}
	if err := json.Unmarshal(queryResponse.BodyBytes, spMetadata); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to FetchSamlSpMetadata: %w", err)
	}

	return spMetadata, nil
}

// SetSamlIdpMetadata will set the Identity Provider metadata needed to configure a SAML connection for an org.
func (o *Client) SetSamlIdpMetadata(params models.SamlIdpMetadata) (bool, error) {
	urlPostfix := "saml_idp_metadata"

	// check that provider is valid
	valid_providers := []string{"Google", "Rippling", "OneLogin", "JumpCloud", "Okta", "Azure", "Duo", "Generic"}
	is_valid_provider := false
	for _, provider := range valid_providers {
		if params.Provider == provider {
			is_valid_provider = true
		}
	}
	if !is_valid_provider {
		return false, fmt.Errorf("Error on setting SAML IDP Metadata for org: provider must be one of %v", valid_providers)
	}

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on setting SAML IDP Metadata for org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on setting SAML IDP Metadata for org: %w", err)
	}

	return true, nil
}

// SamlGoLive will set the SAML connection for an org to be live.
func (o *Client) SamlGoLive(orgID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("saml_idp_metadata/go_live/%s", orgID)

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on setting SAML connection to live for org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on setting SAML connection to live for org: %w", err)
	}

	return true, nil
}

// DeleteSamlConnection will delete the SAML connection for an org.
func (o *Client) DeleteSamlConnection(orgID uuid.UUID) (bool, error) {
	urlPostfix := fmt.Sprintf("saml_idp_metadata/%s", orgID)

	queryResponse, err := o.queryHelper.Delete(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on deleting SAML connection for org: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on deleting SAML connection for org: %w", err)
	}

	return true, nil
}

// public methods for managing API Keys

func (o *Client) FetchAPIKey(apiKeyID string) (*models.APIKeyFull, error) {
	urlPostfix := fmt.Sprintf("end_user_api_keys/%s", apiKeyID)

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, nil)
	if err != nil {
		return nil, fmt.Errorf("Error on fetching an API key: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on fetching an API key: %w", err)
	}

	apiKey := &models.APIKeyFull{}
	if err := json.Unmarshal(queryResponse.BodyBytes, apiKey); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to APIKeyFull: %w", err)
	}

	return apiKey, nil
}

func (o *Client) CreateAPIKey(params models.APIKeyCreateParams) (*models.APIKeyNew, error) {
	urlPostfix := "end_user_api_keys"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating an API key: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating an API key: %w", err)
	}

	apiKey := &models.APIKeyNew{}
	if err := json.Unmarshal(queryResponse.BodyBytes, apiKey); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to APIKeyNew: %w", err)
	}

	return apiKey, nil
}

func (o *Client) UpdateAPIKey(apiKeyID string, params models.APIKeyUpdateParams) (bool, error) {
	urlPostfix := fmt.Sprintf("end_user_api_keys/%s", apiKeyID)

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return false, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Patch(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return false, fmt.Errorf("Error on updating an API key: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on updating an API key: %w", err)
	}

	return true, nil
}

func (o *Client) DeleteAPIKey(apiKeyID string) (bool, error) {
	urlPostfix := fmt.Sprintf("end_user_api_keys/%s", apiKeyID)

	queryResponse, err := o.queryHelper.Delete(o.integrationAPIKey, urlPostfix, nil, nil)
	if err != nil {
		return false, fmt.Errorf("Error on deleting an API key: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return false, fmt.Errorf("Error on deleting an API key: %w", err)
	}

	return true, nil
}

func (o *Client) FetchCurrentAPIKeys(params models.APIKeysQueryParams) (*models.APIKeyResultPage, error) {
	urlPostfix := "end_user_api_keys"

	// assemble the parameters

	queryParams := url.Values{}

	if params.PageNumber != nil {
		queryParams.Add("page_number", strconv.Itoa(*params.PageNumber))
	}
	if params.PageSize != nil {
		queryParams.Add("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.UserID != nil {
		queryParams.Add("user_id", params.UserID.String())
	}
	if params.UserEmail != nil {
		queryParams.Add("user_email", *params.UserEmail)
	}
	if params.OrgID != nil {
		queryParams.Add("org_id", params.OrgID.String())
	}

	// make the request

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on querying API keys: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on querying API keys: %w", err)
	}

	fmt.Println(string(queryResponse.BodyText))

	apiKeys := &models.APIKeyResultPage{}
	if err := json.Unmarshal(queryResponse.BodyBytes, apiKeys); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to APIKeyResultPage: %w", err)
	}

	return apiKeys, nil
}

func (o *Client) FetchArchivedAPIKeys(params models.APIKeysQueryParams) (*models.APIKeyResultPage, error) {
	urlPostfix := "end_user_api_keys/archived"

	// assemble the parameters

	queryParams := url.Values{}

	if params.PageNumber != nil {
		queryParams.Add("page_number", strconv.Itoa(*params.PageNumber))
	}
	if params.PageSize != nil {
		queryParams.Add("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.UserID != nil {
		queryParams.Add("user_id", params.UserID.String())
	}
	if params.UserEmail != nil {
		queryParams.Add("user_email", *params.UserEmail)
	}
	if params.OrgID != nil {
		queryParams.Add("org_id", params.OrgID.String())
	}

	// make the request

	queryResponse, err := o.queryHelper.Get(o.integrationAPIKey, urlPostfix, queryParams)
	if err != nil {
		return nil, fmt.Errorf("Error on querying archived API keys: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on querying archived API keys: %w", err)
	}

	apiKeys := &models.APIKeyResultPage{}
	if err := json.Unmarshal(queryResponse.BodyBytes, apiKeys); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to APIKeyResultPage: %w", err)
	}

	return apiKeys, nil
}

func (o *Client) ValidatePersonalAPIKey(apiKeyToken string) (*models.PersonalAPIKeyValidation, error) {
	apiKeyValidate, err := o.ValidateAPIKey(apiKeyToken)
	if err != nil {
		return nil, err
	}
	if apiKeyValidate.Org != nil || apiKeyValidate.User == nil {
		return nil, fmt.Errorf("not a personal API Key")
	}
	return &models.PersonalAPIKeyValidation{
		User:     *apiKeyValidate.User,
		Metadata: apiKeyValidate.Metadata,
	}, nil
}

func (o *Client) ValidateOrgAPIKey(apiKeyToken string) (*models.OrgAPIKeyValidation, error) {
	apiKeyValidate, err := o.ValidateAPIKey(apiKeyToken)
	if err != nil {
		return nil, err
	}
	if apiKeyValidate.Org == nil {
		return nil, fmt.Errorf("not an org API Key")
	}
	return &models.OrgAPIKeyValidation{
		Org:       *apiKeyValidate.Org,
		Metadata:  apiKeyValidate.Metadata,
		User:      apiKeyValidate.User,
		UserInOrg: apiKeyValidate.UserInOrg,
	}, nil
}

func (o *Client) ValidateAPIKey(apiKeyToken string) (*models.APIKeyValidation, error) {
	urlPostfix := "end_user_api_keys/validate"

	// assemble the parameters

	type ValidteAPIKey struct {
		APIKeyToken string `json:"api_key_token"`
	}

	bodyParams := ValidteAPIKey{
		APIKeyToken: apiKeyToken,
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on validating an API Key: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		if queryResponse.StatusCode == 429 {
			rateLimitError := &models.ApiKeyRateLimitError{}
			err = json.Unmarshal(queryResponse.BodyBytes, rateLimitError)
			if err != nil {
				return nil, fmt.Errorf("Error on unmarshalling bytes to ApiKeyRateLimitError: %w", err)
			} else {
				return nil, rateLimitError
			}
		}
		return nil, fmt.Errorf("Error on validating an API Key: %w", err)
	}

	apiKeyValidate := &models.APIKeyValidation{}
	if err := json.Unmarshal(queryResponse.BodyBytes, apiKeyValidate); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to APIKeyValidation: %w", err)
	}

	return apiKeyValidate, nil
}

// public methods for misc functionality

// CreateAccessToken creates an access token.
func (o *Client) CreateAccessToken(userID uuid.UUID, durationInMinutes int, createAccessTokenOptions ...models.CreateAccessTokenOptions) (*models.AccessToken, error) {
	urlPostfix := "access_token"

	// assemble body params

	type CreateAccessToken struct {
		UserID            uuid.UUID  `json:"user_id"`
		DurationInMinutes int        `json:"duration_in_minutes"`
		ActiveOrgId       *uuid.UUID `json:"active_org_id,omitempty"`
	}

	bodyParams := CreateAccessToken{
		UserID:            userID,
		DurationInMinutes: durationInMinutes,
	}

	if len(createAccessTokenOptions) == 1 {
		bodyParams.ActiveOrgId = createAccessTokenOptions[0].ActiveOrgId
	}

	bodyJSON, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	// make the request

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating access token: %w", err)
	}

	if queryResponse.StatusCode != 200 {
		switch statusCode := queryResponse.StatusCode; statusCode {
		case 401:
			return nil, fmt.Errorf("API Key is incorrect")
		case 400:
			return nil, fmt.Errorf("Bad request: %s", queryResponse.BodyText)
		case 403:
			return nil, fmt.Errorf("User not found")
		case 404:
			return nil, fmt.Errorf("Access token creation not enabled")
		default:
			return nil, fmt.Errorf("Unknown error when creating access token. Status code: %s. Body: %s", strconv.Itoa(queryResponse.StatusCode), queryResponse.BodyText)
		}
	}

	accessToken := &models.AccessToken{}
	if err := json.Unmarshal(queryResponse.BodyBytes, accessToken); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to AccessToken: %w", err)
	}

	return accessToken, nil
}

// CreateMagicLink will create (but not send) a link to let a user sign in without a password.
func (o *Client) CreateMagicLink(params models.CreateMagicLinkParams) (*models.CreateMagicLinkResponse, error) {
	urlPostfix := "magic_link"

	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error on marshalling body params: %w", err)
	}

	queryResponse, err := o.queryHelper.Post(o.integrationAPIKey, urlPostfix, nil, bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error on creating magic link: %w", err)
	}

	if err := o.returnErrorMessageIfNotOk(queryResponse); err != nil {
		return nil, fmt.Errorf("Error on creating magic link: %w", err)
	}

	magicLink := &models.CreateMagicLinkResponse{}
	if err := json.Unmarshal(queryResponse.BodyBytes, magicLink); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to CreateMagicLinkResponse: %w", err)
	}

	return magicLink, nil
}

// public methods around authorization

// GetUser will get a user from a JWT token. From there you get orgs the user is in, and validate the user's
// permissions or roles. See the UserFromToken type for more info.
func (o *Client) GetUser(authHeader string) (*models.UserFromToken, error) {
	accessToken, err := o.validationHelper.ExtractTokenFromAuthorizationHeader(authHeader)
	if err != nil {
		return nil, fmt.Errorf("Error on extracting token from authorization header: %w", err)
	}

	user, err := o.validationHelper.ValidateAccessTokenAndGetUser(accessToken, o.tokenVerificationMetadata)
	if err != nil {
		return nil, fmt.Errorf("Error on validating access token and getting user: %w", err)
	}

	return user, nil
}

// private method to handle errors

func (o *Client) returnErrorMessageIfNotOk(queryResponse *helpers.QueryResponse) error {
	if queryResponse.StatusCode != 200 {
		switch statusCode := queryResponse.StatusCode; statusCode {
		case 401:
			return fmt.Errorf("API Key is incorrect")
		case 400:
			return fmt.Errorf("Bad request: %s", queryResponse.BodyText)
		case 404:
			return fmt.Errorf("API not found")
		case 426:
			return fmt.Errorf("Cannot use organizations unless B2B support is enabled--enable it in your PropelAuth dashboard")
		case 429:
			return fmt.Errorf("Your app is making too many requests, too quickly")
		default:
			return fmt.Errorf("Unknown error when performing operation. Status code: %s. Body: %s", strconv.Itoa(queryResponse.StatusCode), queryResponse.BodyText)
		}
	}

	return nil
}
