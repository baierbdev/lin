package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var fs http.FileSystem = http.Dir("./data")
func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	status := c.Param("status")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error})
		return
	}
	id := uuid.New();

	dst := filepath.Join("./data/", filepath.Base(id.String()+"-"+status+"-"+file.Filename))
	c.SaveUploadedFile(file, dst)
	
	c.String(http.StatusOK, id.String()+"-"+status+"-"+file.Filename))
}

func DownloadDocument(c *gin.Context) {
	c.FileFromFS(c.Param("name"), fs)
}
