package payload

type CreateLeadRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
}

type CreateLeadResponse struct {
	LeadID      string `json:"lead_id"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}
