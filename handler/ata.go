package handler

import (
	"lin/models"
	"lin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AtaHandler struct {
	ataService service.AtaService
}

func NewAtaHandler(svc *service.AtaService) *AtaHandler {
	return &AtaHandler{
		ataService: *svc,
	}
}
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
func (h *AtaHandler) DownloadAta(c *gin.Context) {
	filename := c.Param("name")
	filePath := h.ataService.GetAta(filename)
	c.File(filePath)
}
func (h *AtaHandler) DeleteAta(c *gin.Context) {
	filename := c.Param("name")

	if err := h.ataService.DeleteAta(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
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
