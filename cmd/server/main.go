package main

// Esse pacote é responsável por gerir os documentos (atas, contratos, aditivos
// e notas ficais) do gestorcpm, e atualmente vem sendo expandido para embarcar
// o portal de compras públicas (pncp).

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"lin/handler"
	"lin/service"

	"github.com/gin-gonic/gin"
)

// corsMiddleware retorna um middleware Gin que configura cabeçalhos CORS
// permitindo requisições de qualquer origem (Access-Control-Allow-Origin: *).
// Também responde automaticamente a requisições OPTIONS (preflight) com 204 No Content.
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

// main é o ponto de entrada da aplicação. Inicializa os serviços e handlers
// para gestão de documentos (notas fiscais, contratos, atas, aditivos),
// registra as rotas no roteador Gin e inicia o servidor HTTP na porta 8080.
func main() {
	if urlPncp == "" {
		log.Fatal("variável de ambiente PNCP_URL não está definida")
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
		log.Fatalf("servidor falhou ao iniciar na porta: %s", err.Error())
		os.Exit(1)
	}
}
