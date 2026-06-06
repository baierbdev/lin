package handler

import (
	"lin/models"
	"lin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AtaHandler gerencia as requisições HTTP relacionadas a atas de registro de preço,
// incluindo upload, download e exclusão de arquivos, além da consulta ao PNCP.
type AtaHandler struct {
	ataService service.AtaService
}

// NewAtaHandler cria um novo AtaHandler com o serviço de atas fornecido.
func NewAtaHandler(svc *service.AtaService) *AtaHandler {
	return &AtaHandler{
		ataService: *svc,
	}
}
// UploadFile gerencia o upload de arquivos de ata (POST /atas).
// Aceita um arquivo multipart e o campo ata_id (UUID válido obrigatório).
// O arquivo é salvo com o prefixo do ata_id seguido do nome original do arquivo.
func (h *AtaHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ataId := c.PostForm("ata_id")
	if _, err := uuid.Parse(ataId); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "ata_id deve ser um UUID válido"})
		return
	}

	if err := h.ataService.EnsureAtaDataDir(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	outputName, err := h.ataService.SaveFile(file, ataId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.String(http.StatusOK, outputName)
}
// DownloadAta serve o arquivo de ata solicitado para download (GET /atas/:name).
// O parâmetro :name na URL corresponde ao nome do arquivo armazenado no diretório de dados.
func (h *AtaHandler) DownloadAta(c *gin.Context) {
	filename := c.Param("name")
	filePath := h.ataService.GetAta(filename)
	c.File(filePath)
}
// DeleteAta remove um arquivo de ata do armazenamento (DELETE /atas/:name).
// Retorna 204 No Content em caso de sucesso ou 500 Internal Server Error em caso de falha.
func (h *AtaHandler) DeleteAta(c *gin.Context) {
	filename := c.Param("name")

	if err := h.ataService.DeleteAta(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
// LoadAtaPncp consulta informações de uma ata no Portal Nacional de Contratações Públicas
// (GET /atas/pncp/:cnpj/:year/:sequencialCompra/:sequencialAta).
// Retorna os dados da ata em JSON ou 404 se não encontrada no PNCP.
func (h *AtaHandler) LoadAtaPncp(c *gin.Context) {
	cnpj := c.Param("cnpj")
	year := c.Param("year")
	sequencialCompra := c.Param("sequencialCompra")
	sequencialAta := c.Param("sequencialAta")

	data, err := h.ataService.GetAtaInfoPncp(cnpj, year, sequencialCompra, sequencialAta)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
