package response

type CreateLeadResponse struct {
	LeadID      string `json:"lead_id"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}
