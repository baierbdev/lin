package handler

import (
	"lin/models"
	"lin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContratoHandler gerencia as requisições HTTP relacionadas a contratos,
// incluindo upload, download e exclusão de arquivos, além da consulta ao PNCP.
type ContratoHandler struct {
	contratoService service.ContratoService
}

// NewContratoHandler cria um novo ContratoHandler com o serviço de contratos fornecido.
func NewContratoHandler(svc *service.ContratoService) *ContratoHandler {
	return &ContratoHandler{
		contratoService: *svc,
	}
}
// UploadFile gerencia o upload de arquivos de contrato (POST /contratos).
// Aceita um arquivo multipart e o campo contrato_id (UUID válido obrigatório).
// O arquivo é salvo com o prefixo do contrato_id seguido do nome original do arquivo.
func (h *ContratoHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	contratoId := c.PostForm("contrato_id")
	if _, err := uuid.Parse(contratoId); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "contrato_id deve ser um UUID válido"})
		return
	}

	if err := h.contratoService.EnsureContratoDataDir(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	outputName, err := h.contratoService.SaveFile(file, contratoId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.String(http.StatusOK, outputName)
}
// DownloadContrato serve o arquivo de contrato solicitado para download (GET /contratos/:name).
// O parâmetro :name na URL corresponde ao nome do arquivo armazenado no diretório de dados.
func (h *ContratoHandler) DownloadContrato(c *gin.Context) {
	filename := c.Param("name")
	filePath := h.contratoService.GetContrato(filename)
	c.File(filePath)
}
// DeleteContrato remove um arquivo de contrato do armazenamento (DELETE /contratos/:name).
// Retorna 204 No Content em caso de sucesso ou 500 Internal Server Error em caso de falha.
func (h *ContratoHandler) DeleteContrato(c *gin.Context) {
	filename := c.Param("name")

	if err := h.contratoService.Deletecontrato(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

// LoadContratoPncp consulta informações de um contrato no Portal Nacional de Contratações Públicas
// (GET /contratos/pncp/:cnpj/:ano/:sequencialContrato).
// Retorna os dados do contrato em JSON ou erro se não encontrado no PNCP.
func (h *ContratoHandler) LoadContratoPncp(c *gin.Context) {
	cnpj := c.Param("cnpj")
	ano := c.Param("ano")
	sequencial := c.Param("sequencialContrato")

	resp, err := h.contratoService.GetContratoPncp(cnpj, ano, sequencial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
