package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lin/internal/models"
)

type DocumentService struct {
	dataDir string
}

func NewDocumentService(dataDir string) *DocumentService {
	return &DocumentService{
		dataDir: dataDir,
	}
}

func (s *DocumentService) EnsureDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

func (s *DocumentService) SaveFile(fileHeader *multipart.FileHeader, notaID, status string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputName := fmt.Sprintf("%s-%s-%s", notaID, status, safeFilename)
	dst := filepath.Join(s.dataDir, outputName)

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return outputName, nil
}

func (s *DocumentService) ListByNotaID(notaID string) ([]models.ListedDocument, error) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.ListedDocument{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	prefix := notaID + "-"
	documents := make([]models.ListedDocument, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		documents = append(documents, models.ListedDocument{
			Name:   name,
			Status: extractStatusFromFilename(name, notaID),
			URL:    fmt.Sprintf("/retrieve/%s", name),
		})
	}

	sort.Slice(documents, func(i, j int) bool {
		return documents[i].Name > documents[j].Name
	})

	return documents, nil
}

func (s *DocumentService) GetFilePath(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

func extractStatusFromFilename(fileName, notaID string) string {
	prefix := notaID + "-"
	if !strings.HasPrefix(fileName, prefix) {
		return ""
	}

	remainder := strings.TrimPrefix(fileName, prefix)
	parts := strings.SplitN(remainder, "-", 2)
	if len(parts) < 2 || parts[0] == "" {
		return ""
	}

	return parts[0]
}
