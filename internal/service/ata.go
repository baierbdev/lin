package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"io"
	"path/filepath"
)


type AtaService struct {
	dataDir string
}
func NewAtaService(dataDir string) *AtaService {
	return &AtaService{
		dataDir: dataDir,
	}
}
func (s *NotaService) EnsureAtaDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

func (s *AtaService) SaveFile(fileHeader *multipart.FileHeader, ataId string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s", ataId, safeFilename )
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
