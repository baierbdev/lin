package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type AditivoService struct {
	dataDir string
}

func NewAditivoService(dataDir string) *AditivoService {
	return &AditivoService{
		dataDir: dataDir,
	}
}

func (s *AditivoService) EnsureAditivoDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

func (s *AditivoService) SaveFile(fileHeader *multipart.FileHeader, date string, tipo, contratoId string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer src.Close()

	safeName := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s-%s-%s", contratoId, tipo, date, safeName)
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

func (s *AditivoService) GetAditivo(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

func (s *AditivoService) DeleteAditivo(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("falha ao remover arquivo: %w", err)
	}
	return nil
}
