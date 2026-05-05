package handler

import (
	"lin/internal/models"
	"lin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AtaHandler struct {
	ataService service.AtaService
}
func NewAtaHandler(svc *service.AtaService) *AtaHandler {
	return  &AtaHandler{
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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "ata_id must be a valid UUID"})
		return
	}

	outputName, err := h.ataService.SaveFile(file, ataId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	c.JSON(http.StatusOK, outputName)
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
