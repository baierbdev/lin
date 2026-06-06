package models

// ErrorResponse é o envelope padrão para respostas de erro da API,
// contendo uma mensagem descritiva no campo "error".
type ErrorResponse struct {
	Error string `json:"error"`
}
