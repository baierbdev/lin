package service

import (
	"encoding/json"
	"fmt"
	"io"
	"lin/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ContratoService struct {
	dataDir string
	client  http.Client
	urlPncp string
}

func NewContratoService(dataDir string, urlPncp string, client http.Client) *ContratoService {
	return &ContratoService{
		dataDir: dataDir,
		client:  client,
		urlPncp: urlPncp,
	}
}

func (s *ContratoService) EnsureContratoDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

func (s *ContratoService) SaveFile(fileHeader *multipart.FileHeader, contratoId string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s", contratoId, safeFilename)
	dst := filepath.Join(s.dataDir, outputFilename)

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo de destino: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	return outputFilename, nil
}

func (s *ContratoService) GetContrato(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

func (s *ContratoService) Deletecontrato(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("falha ao remover arquivo: %w", err)
	}
	return nil
}

func (s *ContratoService) GetContratoPncp(cnpj string, ano string, sequencialContrato string) (*models.ContratoPncp, error) {
	urlReq := fmt.Sprintf("%s/v1/orgaos/%s/contratos/%s/%s",
		s.urlPncp, cnpj, ano, sequencialContrato,
	)
	res, err := s.client.Get(urlReq)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("contrato não encontrado no PNCP (404)")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do PNCP retornou status inesperado: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao processar a resposta: %w", err)
	}

	data := &models.ContratoPncp{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, fmt.Errorf("falha ao recuperar contrato do PNCP: %w", err)
	}
	return data, nil
}
