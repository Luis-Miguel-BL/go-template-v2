package response

type AuthorizationResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}
