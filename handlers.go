package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var fs http.FileSystem = http.Dir("./data")

func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := c.Param("status")
	notaID := c.PostForm("nota_id")
	if notaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nota_id is required"})
		return
	}

	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nota_id must be a valid UUID"})
		return
	}

	if err := os.MkdirAll("./data", 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	safeFilename := filepath.Base(file.Filename)
	outputName := notaID + "-" + status + "-" + safeFilename
	dst := filepath.Join("./data", outputName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, outputName)
}

func DownloadDocument(c *gin.Context) {
	c.FileFromFS(c.Param("name"), fs)
}
