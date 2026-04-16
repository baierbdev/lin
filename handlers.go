package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var fs http.FileSystem = http.Dir("./data")

type ListedDocument struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

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

func ListDocumentsByNota(c *gin.Context) {
	notaID := c.Param("nota_id")
	if _, err := uuid.Parse(notaID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nota_id must be a valid UUID"})
		return
	}

	entries, err := os.ReadDir("./data")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{
				"nota_id":   notaID,
				"count":     0,
				"documents": []ListedDocument{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	prefix := notaID + "-"
	documents := make([]ListedDocument, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		documents = append(documents, ListedDocument{
			Name:   name,
			Status: extractStatusFromFilename(name, notaID),
			URL:    "/retrieve/" + url.PathEscape(name),
		})
	}

	sort.Slice(documents, func(i, j int) bool {
		return documents[i].Name > documents[j].Name
	})

	c.JSON(http.StatusOK, gin.H{
		"nota_id":   notaID,
		"count":     len(documents),
		"documents": documents,
	})
}

func extractStatusFromFilename(fileName, notaID string) string {
	prefix := notaID + "-"
	if !strings.HasPrefix(fileName, prefix) {
		return ""
	}

	remainder := strings.TrimPrefix(fileName, prefix)
	parts := strings.SplitN(remainder, "-", 2)
	if len(parts) < 2 || parts[0] == "" {
		return ""
	}

	return parts[0]
}
