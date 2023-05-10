package models

// return types

// CreateMagicLinkResponse has one field, URL, which is the magic link to sign someone in automatically.
type CreateMagicLinkResponse struct {
	Url string `json:"url"`
}

// post types

// CreateMagicLinkParams is the information needed to create a magic link to sign someone in automatically.
type CreateMagicLinkParams struct {
	Email                         string  `json:"email"`
	RedirectURL                   *string `json:"redirect_url,omitempty"`
	ExpiresInHours                *int    `json:"expires_in_hours,omitempty"`
	CreateNewUserIfOneDoesntExist *bool   `json:"create_new_user_if_one_doesnt_exist,omitempty"`
}
