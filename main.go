package main

// Esse pacote é responsável por gerir os documentos (atas, contratos, aditivos
// e notas ficais) do gestorcpm, e atualmente vem sendo expandido para embarcar
// o portal de compras públicas (pncp).

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"lin/internal/handler"
	"lin/internal/service"

	"github.com/gin-gonic/gin"
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

var rootDir = "./data"
var urlPncp = os.Getenv("PNCP_URL")

func main() {
	if urlPncp == "" {
		log.Fatal("PNCP_URL environment variable is not set")
	}

	router := gin.Default()
	router.Use(corsMiddleware())

	client := &http.Client{}

	notaService := service.NewNotaService(filepath.Join(rootDir, "notas"))
	contratoService := service.NewContratoService(filepath.Join(rootDir, "contratos"), urlPncp, *client)
	atasServices := service.NewAtaService(filepath.Join(rootDir, "atas"), urlPncp, *client)
	aditivoServices := service.NewAditivoService(filepath.Join(rootDir, "aditivos"))

	notaHandler := handler.NewNotaHandler(notaService)
	contratoHandler := handler.NewContratoHandler(contratoService)
	ataHandler := handler.NewAtaHandler(atasServices)
	aditivoHandler := handler.NewAditivoService(aditivoServices)

	router.POST("/notas/upload/:status", notaHandler.UploadNota)
	router.GET("/notas/retrieve/:name", notaHandler.DownloadNota)
	router.GET("/notas/list/:nota_id", notaHandler.ListNotasByNota)

	router.GET("/contratos/pncp/:cnpj/:ano/:sequencialContrato", contratoHandler.LoadContratoPncp)
	router.POST("/contratos", contratoHandler.UploadFile)
	router.GET("/contratos/:name", contratoHandler.DownloadContrato)
	router.DELETE("/contratos/:name", contratoHandler.DeleteContrato)

	router.GET("/atas/pncp/:cnpj/:year/:sequencialCompra/:sequencialAta", ataHandler.LoadAtaPncp)
	router.POST("/atas", ataHandler.UploadFile)
	router.GET("/atas/:name", ataHandler.DownloadAta)
	router.DELETE("/atas/:name", ataHandler.DeleteAta)

	router.POST("/aditivos", aditivoHandler.UploadFile)
	router.GET("/aditivos/:name", aditivoHandler.DownloadAditivo)
	router.DELETE("/aditivos/:name", aditivoHandler.DeleteAditivo)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start in port: %s", err.Error())
		os.Exit(1)
	}
}
