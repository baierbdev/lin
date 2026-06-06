package models

// ListedNota representa um arquivo de nota fiscal na listagem,
// com nome, status extraído do nome do arquivo e URL relativa para download.
type ListedNota struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

// NotaListResponse é o envelope de resposta para a listagem de notas fiscais
// por nota_id, contendo o identificador da nota, a quantidade de arquivos
// e a lista de arquivos encontrados.
type NotaListResponse struct {
	NotaID string       `json:"nota_id"`
	Count  int          `json:"count"`
	Notas  []ListedNota `json:"notas"`
}
