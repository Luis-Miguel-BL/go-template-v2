package request

type CreateLeadRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
}
