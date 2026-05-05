package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"io"
	"path/filepath"
)

type ContratoService struct {
	dataDir string 
}
func NewContratoService(dataDir string) *ContratoService {
	return &ContratoService{
		dataDir: dataDir,
	}		
}

func (s *NotaService) EnsureContratoDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}
func (s *ContratoService) SaveFile(fileHeader *multipart.FileHeader, contratoId string) (string, error)  {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open a file: %w", err)
	}
	defer src.Close()
	

	safeFilename := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s", contratoId, safeFilename)
	dst := filepath.Join(s.dataDir, outputFilename) 

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	} 
	out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return outputFilename, nil
}
func (s *ContratoService) GetContrato(filename string) (string) {
	return  filepath.Join(s.dataDir, filename) 
}
func (s *ContratoService) Deletecontrato(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil	
} 
