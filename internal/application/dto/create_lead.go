package dto

type CreateLeadInput struct {
	Name           string
	Email          string
	Phone          string
	DocumentNumber string
}

type CreateLeadOutput struct {
	LeadID      string
	AccessToken AccessToken
}
