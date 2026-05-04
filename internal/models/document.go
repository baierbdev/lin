package models

type ListedDocument struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type DocumentListResponse struct {
	NotaID    string           `json:"nota_id"`
	Count     int              `json:"count"`
	Documents []ListedDocument `json:"documents"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
