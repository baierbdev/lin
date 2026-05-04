package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lin/internal/handler"
	"lin/internal/service"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())

	docService := service.NewDocumentService("./data")
	docHandler := handler.NewDocumentHandler(docService)

	router.POST("/upload/:status", docHandler.UploadDocument)
	router.GET("/retrieve/:name", docHandler.DownloadDocument)
	router.GET("/list/:nota_id", docHandler.ListDocumentsByNota)

	router.Run(":8080")
}
