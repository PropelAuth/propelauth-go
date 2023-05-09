package models

// return types

type CreateMagicLinkResponse struct {
	Url string `json:"url"`
}

// post types

type CreateMagicLinkParams struct {
	Email                         string  `json:"email"`
	RedirectURL                   *string `json:"redirect_url,omitempty"`
	ExpiresInHours                *int    `json:"expires_in_hours,omitempty"`
	CreateNewUserIfOneDoesntExist *bool   `json:"create_new_user_if_one_doesnt_exist,omitempty"`
}
