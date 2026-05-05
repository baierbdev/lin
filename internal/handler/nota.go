package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lin/internal/models"
	"lin/internal/service"
)

type NotaHandler struct {
	service *service.NotaService
}

func NewNotaHandler(svc *service.NotaService) *NotaHandler {
	return &NotaHandler{
		service: svc,
	}
}

func (h *NotaHandler) UploadNota(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	status := c.Param("status")
	notaID := c.PostForm("nota_id")
	if notaID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id is required"})
		return
	}

	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id must be a valid UUID"})
		return
	}

	if err := h.service.EnsureDataDir(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	outputName, err := h.service.SaveFile(file, notaID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.String(http.StatusOK, outputName)
}

func (h *NotaHandler) DownloadNota(c *gin.Context) {
	fileName := c.Param("name")
	filePath := h.service.GetFilePath(fileName)
	c.File(filePath)
}

func (h *NotaHandler) ListNotasByNota(c *gin.Context) {
	notaID := c.Param("nota_id")
	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id must be a valid UUID"})
		return
	}

	notas, err := h.service.ListByNotaID(notaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NotaListResponse{
		NotaID:  notaID,
		Count:   len(notas),
		Notas:   notas,
	})
}
