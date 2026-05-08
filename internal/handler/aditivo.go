package handler

import (
	"lin/internal/models"
	"lin/internal/service"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AditivoHandler struct {
	aditivoService *service.AditivoService
}

func NewAditivoService(srv *service.AditivoService) *AditivoHandler {
	return &AditivoHandler{
		aditivoService: srv,
	}
}
func (h *AditivoHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	tipo := c.PostForm("tipo")
	if tipo == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "tipo is required"})
		return
	}

	contratoID := c.PostForm("contrato_id")
	if _, err := uuid.Parse(contratoID); err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "contrato_id must be a valid UUID"})
		return
	}
	if err := h.aditivoService.EnsureAditivoDataDir(); err != nil {

		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	date := c.PostForm("data")
	if _, err := time.Parse(time.RFC3339, date); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "date must be a valid date"})
		return

	}

	outputName, err := h.aditivoService.SaveFile(file, date, tipo, contratoID)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.String(http.StatusOK, outputName)
}

func (h *AditivoHandler) DownloadAditivo(c *gin.Context) {
	filename := c.Param("name")
	filepath := h.aditivoService.GetAditivo(filename)
	c.File(filepath)
}
func (h *AditivoHandler) DeleteAditivo(c *gin.Context) {
	filename := c.Param("name")

	if err := h.aditivoService.DeleteAditivo(filename); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
