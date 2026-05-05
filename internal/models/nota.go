package models

type ListedNota struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type NotaListResponse struct {
	NotaID    string         `json:"nota_id"`
	Count     int            `json:"count"`
	Notas     []ListedNota   `json:"notas"`
}

