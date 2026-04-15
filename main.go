package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.POST("/upload", UploadDocument)
	router.GET("/retrieve/:name", DownloadDocument)

	router.Run(":8080")
}
