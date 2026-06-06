package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lin/models"
	"lin/service"
)

// NotaHandler gerencia as requisições HTTP relacionadas a notas fiscais,
// incluindo upload, download e listagem de arquivos por nota.
type NotaHandler struct {
	service *service.NotaService
}

// NewNotaHandler cria um novo NotaHandler com o serviço de notas fiscais fornecido.
func NewNotaHandler(svc *service.NotaService) *NotaHandler {
	return &NotaHandler{
		service: svc,
	}
}

// UploadNota gerencia o upload de arquivos de nota fiscal (POST /notas/upload/:status).
// Aceita múltiplos arquivos multipart e o campo nota_id (UUID válido obrigatório).
// O status é extraído da URL e aplicado apenas ao primeiro arquivo; os demais
// são salvos sem status. Retorna a URL do último arquivo processado.
func (h *NotaHandler) UploadNota(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	status := c.Param("status")
	notaID := c.PostForm("nota_id")
	if notaID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id é obrigatório"})
		return
	}

	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id deve ser um UUID válido"})
		return
	}

	if err := h.service.EnsureDataDir(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var fileUrl string

	for key, file := range form.File["files"] {
		if key != 0 {
			status = ""
		}
		outputName, err := h.service.SaveFile(file, notaID, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		fileUrl = outputName
	}

	c.String(http.StatusOK, fileUrl)
}

// DownloadNota serve o arquivo de nota fiscal solicitado para download (GET /notas/retrieve/:name).
// O parâmetro :name na URL corresponde ao nome do arquivo armazenado no diretório de dados.
func (h *NotaHandler) DownloadNota(c *gin.Context) {
	fileName := c.Param("name")
	filePath := h.service.GetFilePath(fileName)
	c.File(filePath)
}

// ListNotasByNota lista todos os arquivos de nota fiscal associados a um determinado
// nota_id (GET /notas/list/:nota_id). Retorna um JSON com o ID da nota, a contagem
// de arquivos e a lista de arquivos com nome, status extraído do nome e URL de download.
func (h *NotaHandler) ListNotasByNota(c *gin.Context) {
	notaID := c.Param("nota_id")
	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id deve ser um UUID válido"})
		return
	}

	notas, err := h.service.ListByNotaID(notaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NotaListResponse{
		NotaID: notaID,
		Count:  len(notas),
		Notas:  notas,
	})
}
