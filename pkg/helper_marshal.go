package client

import (
	"encoding/json"
	"fmt"
)

type MarshalHelper struct{}

func (o *MarshalHelper) GetUserFromBytes(bytes []byte) (*UserID, error) {
	var user UserID
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserID`: %v", err)
	}
	return &user, nil
}

func (o *MarshalHelper) GetUserMetadataFromBytes(bytes []byte) (*UserMetadata, error) {
	var user UserMetadata
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserMetadata`: %v", err)
	}
	return &user, nil
}

func (o *MarshalHelper) GetUserListFromBytes(bytes []byte) (*UserList, error) {
	var users UserList
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to UserLIst`: %v", err)
	}
	return &users, nil
}

func (o *MarshalHelper) GetOrgMetadataFromBytes(bytes []byte) (*OrgMetadata, error) {
	var org OrgMetadata
	if err := json.Unmarshal(bytes, &org); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to Org`: %v", err)
	}
	return &org, nil
}

func (o *MarshalHelper) GetOrgListFromBytes(bytes []byte) (*OrgList, error) {
	var orgs OrgList
	if err := json.Unmarshal(bytes, &orgs); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to OrgList`: %v", err)
	}
	return &orgs, nil
}

func (o *MarshalHelper) GetAccessTokenFromBytes(bytes []byte) (*AccessToken, error) {
	var accessToken AccessToken
	if err := json.Unmarshal(bytes, &accessToken); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to AccessToken`: %v", err)
	}
	return &accessToken, nil
}

func (o *MarshalHelper) GetMagicLinkResponseFromBytes(bytes []byte) (*CreateMagicLinkResponse, error) {
	var magicLink CreateMagicLinkResponse
	if err := json.Unmarshal(bytes, &magicLink); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to MagicLinkResponse`: %v", err)
	}
	return &magicLink, nil
}

func (o *MarshalHelper) GetAuthTokenVerificationMetadataResponseFromBytes(bytes []byte) (*AuthTokenVerificationMetadataResponse, error) {
	var authTokenVerificationMetadataResponse AuthTokenVerificationMetadataResponse
	if err := json.Unmarshal(bytes, &authTokenVerificationMetadataResponse); err != nil {
		return nil, fmt.Errorf("Error on unmarshalling bytes to AuthTokenVerificationMetadataResponse`: %v", err)
	}
	return &authTokenVerificationMetadataResponse, nil
}
