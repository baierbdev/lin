package handler

import (
	"lin/models"
	"lin/service"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AditivoHandler gerencia as requisições HTTP relacionadas a aditivos contratuais,
// incluindo upload, download e exclusão de arquivos de aditivos.
type AditivoHandler struct {
	aditivoService *service.AditivoService
}

// NewAditivoService cria um novo AditivoHandler com o serviço de aditivos fornecido.
func NewAditivoService(srv *service.AditivoService) *AditivoHandler {
	return &AditivoHandler{
		aditivoService: srv,
	}
}
// UploadFile gerencia o upload de arquivos de aditivo (POST /aditivos).
// Aceita um arquivo multipart e os campos: tipo (obrigatório), contrato_id (UUID válido)
// e data (formato RFC3339). O arquivo é salvo com o nome composto por esses metadados.
func (h *AditivoHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	tipo := c.PostForm("tipo")
	if tipo == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "tipo é obrigatório"})
		return
	}

	contratoID := c.PostForm("contrato_id")
	if _, err := uuid.Parse(contratoID); err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "contrato_id deve ser um UUID válido"})
		return
	}
	if err := h.aditivoService.EnsureAditivoDataDir(); err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	date := c.PostForm("data")
	if _, err := time.Parse(time.RFC3339, date); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "data deve ser uma data válida"})
		return

	}

	outputName, err := h.aditivoService.SaveFile(file, date, tipo, contratoID)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.String(http.StatusOK, outputName)
}

// DownloadAditivo serve o arquivo de aditivo solicitado para download (GET /aditivos/:name).
// O parâmetro :name na URL corresponde ao nome do arquivo armazenado no diretório de dados.
func (h *AditivoHandler) DownloadAditivo(c *gin.Context) {
	filename := c.Param("name")
	filepath := h.aditivoService.GetAditivo(filename)
	c.File(filepath)
}
// DeleteAditivo remove um arquivo de aditivo do armazenamento (DELETE /aditivos/:name).
// Retorna 204 No Content em caso de sucesso ou 500 Internal Server Error em caso de falha.
func (h *AditivoHandler) DeleteAditivo(c *gin.Context) {
	filename := c.Param("name")

	if err := h.aditivoService.DeleteAditivo(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
