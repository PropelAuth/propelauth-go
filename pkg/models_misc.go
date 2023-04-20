package client

type CreateMagicLinkParams struct {
	Email                         string `json:"email"`
	RedirectURL                   string `json:"redirect_url"`
	ExpiresInHours                int    `json:"expires_in_hours"`
	CreateNewUserIfOneDoesntExist bool   `json:"create_new_user_if_one_doesnt_exist"`
}

type CreateMagicLinkResponse struct {
	Url string `json:"url"`
}

type AccessTokenResponse struct {
	AccessToken AccessTokenData `json:"access_token"`
}

type AccessTokenData struct {
	AccessToken          string                            `json:"access_token"`
	ExpiresAtSeconds     int64                             `json:"expires_at_seconds"`
	OrgIdToOrgMemberInfo map[string]OrgMemberInfoFromToken `json:"org_id_to_org_member_info"`
	User                 UserMetadata                      `json:"user"`
	ImpersonatorUser     UserID                            `json:"impersonator_user,omitempty"`
}
