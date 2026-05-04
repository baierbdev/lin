package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lin/internal/models"
	"lin/internal/service"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(svc *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: svc,
	}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
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

func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	fileName := c.Param("name")
	filePath := h.service.GetFilePath(fileName)
	c.File(filePath)
}

func (h *DocumentHandler) ListDocumentsByNota(c *gin.Context) {
	notaID := c.Param("nota_id")
	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "nota_id must be a valid UUID"})
		return
	}

	documents, err := h.service.ListByNotaID(notaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.DocumentListResponse{
		NotaID:    notaID,
		Count:     len(documents),
		Documents: documents,
	})
}
