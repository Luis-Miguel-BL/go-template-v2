package exampleapi

type CreateExampleRequestDTO struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type CreateExampleResponseDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}
