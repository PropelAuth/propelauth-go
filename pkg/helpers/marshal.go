package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/propelauth/propelauth-go/pkg/models"
)

type MarshalHelperInterface interface {
	GetUserFromBytes(bytes []byte) (*models.UserID, error)
	GetUserMetadataFromBytes(bytes []byte) (*models.UserMetadata, error)
	GetUserListFromBytes(bytes []byte) (*models.UserList, error)
	GetOrgMetadataFromBytes(bytes []byte) (*models.OrgMetadata, error)
	GetOrgListFromBytes(bytes []byte) (*models.OrgList, error)
	GetAccessTokenFromBytes(bytes []byte) (*models.AccessToken, error)
	GetMagicLinkResponseFromBytes(bytes []byte) (*models.CreateMagicLinkResponse, error)
	GetAuthTokenVerificationMetadataResponseFromBytes(bytes []byte) (*models.AuthTokenVerificationMetadataResponse, error)
}

type MarshalHelper struct{}

// TODO these can be removed and replaced with a generic function that takes a pointer to the type

func (o *MarshalHelper) GetUserFromBytes(bytes []byte) (*models.UserID, error) {
	var user models.UserID
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserID`: %v", err)
	}
	return &user, nil
}

func (o *MarshalHelper) GetUserMetadataFromBytes(bytes []byte) (*models.UserMetadata, error) {
	var user models.UserMetadata
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserMetadata`: %v", err)
	}
	return &user, nil
}

func (o *MarshalHelper) GetUserListFromBytes(bytes []byte) (*models.UserList, error) {
	var users models.UserList
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserLIst`: %v", err)
	}
	return &users, nil
}

func (o *MarshalHelper) GetOrgMetadataFromBytes(bytes []byte) (*models.OrgMetadata, error) {
	var org models.OrgMetadata
	if err := json.Unmarshal(bytes, &org); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to Org`: %v", err)
	}
	return &org, nil
}

func (o *MarshalHelper) GetOrgListFromBytes(bytes []byte) (*models.OrgList, error) {
	var orgs models.OrgList
	if err := json.Unmarshal(bytes, &orgs); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to OrgList`: %v", err)
	}
	return &orgs, nil
}

func (o *MarshalHelper) GetAccessTokenFromBytes(bytes []byte) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	if err := json.Unmarshal(bytes, &accessToken); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to AccessToken`: %v", err)
	}
	return &accessToken, nil
}

func (o *MarshalHelper) GetMagicLinkResponseFromBytes(bytes []byte) (*models.CreateMagicLinkResponse, error) {
	var magicLink models.CreateMagicLinkResponse
	if err := json.Unmarshal(bytes, &magicLink); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to MagicLinkResponse`: %v", err)
	}
	return &magicLink, nil
}

func (o *MarshalHelper) GetAuthTokenVerificationMetadataResponseFromBytes(bytes []byte) (*models.AuthTokenVerificationMetadataResponse, error) {
	var authTokenVerificationMetadataResponse models.AuthTokenVerificationMetadataResponse
	if err := json.Unmarshal(bytes, &authTokenVerificationMetadataResponse); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to AuthTokenVerificationMetadataResponse`: %v", err)
	}
	return &authTokenVerificationMetadataResponse, nil
}
