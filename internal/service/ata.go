package service

import (
	"encoding/json"
	"fmt"
	"io"
	"lin/internal/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type AtaService struct {
	dataDir string
	client  http.Client
	urlPncp string
}

func NewAtaService(dataDir, urlPncp string, client http.Client) *AtaService {
	return &AtaService{
		dataDir: dataDir,
		client:  client,
		urlPncp: urlPncp,
	}
}
func (s *AtaService) EnsureAtaDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

func (s *AtaService) SaveFile(fileHeader *multipart.FileHeader, ataId string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s", ataId, safeFilename)
	dst := filepath.Join(s.dataDir, outputFilename)

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)

	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return outputFilename, nil
}
func (s *AtaService) GetAta(filename string) string {
	return filepath.Join(s.dataDir, filename)
}
func (s *AtaService) DeleteAta(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}

func (s *AtaService) GetAtaInfoPncp(cnpj string, year string, sequencialCompra string, sequencialAta string) (*models.AtaPncp, error) {

	urlReq := fmt.Sprintf("%s/v1/orgaos/%s/compras/%s/%s/atas/%s",
		s.urlPncp, cnpj, year, sequencialCompra, sequencialAta)

	resp, err := s.client.Get(urlReq)
	if err != nil {
		return nil, fmt.Errorf("failed in requisition: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("ata not found in PNCP (404)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PNCP API returned unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing the response: %w", err)
	}

	data := &models.AtaPncp{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, fmt.Errorf("failed to retrieve ata from pncp: %w", err)
	}

	return data, nil
}
