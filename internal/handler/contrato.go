package handler

import (
	"lin/internal/models"
	"lin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContratoHandler struct {
	contratoService service.ContratoService
}
func NewContratoHandler(svc *service.ContratoService) *ContratoHandler {
	return  &ContratoHandler{
		contratoService: *svc,
	}
}
func (h *ContratoHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	
	contratoId := c.PostForm("contrato_id")
	if _, err := uuid.Parse(contratoId); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "contrato_id must be a valid UUID"})
		return
	}

	outputName, err := h.contratoService.SaveFile(file, contratoId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	c.JSON(http.StatusOK, outputName)
} 
func (h *ContratoHandler) DownloadContrato(c *gin.Context) {
	filename := c.Param("name")
	filePath := h.contratoService.GetContrato(filename)
	c.File(filePath)	
}
func (h *ContratoHandler) DeleteContrato(c *gin.Context) {
	filename := c.Param("name")

	if err := h.contratoService.Deletecontrato(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
