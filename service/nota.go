package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lin/models"
)

// NotaService gerencia o armazenamento de arquivos de notas fiscais em disco
// e fornece operações de listagem e recuperação por nota_id.
type NotaService struct {
	dataDir string
}

// NewNotaService cria um novo NotaService que utilizará o diretório de dados especificado.
func NewNotaService(dataDir string) *NotaService {
	return &NotaService{
		dataDir: dataDir,
	}
}

// EnsureDataDir garante que o diretório de dados de notas fiscais exista,
// criando-o com permissões 0755 se necessário.
func (s *NotaService) EnsureDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

// SaveFile salva um arquivo de nota fiscal no diretório de dados com o nome composto
// no formato "{notaID}-{status}-{nomeOriginal}". Retorna o nome do arquivo gerado
// ou um erro em caso de falha na abertura, criação ou cópia do arquivo.
func (s *NotaService) SaveFile(fileHeader *multipart.FileHeader, notaID, status string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputName := fmt.Sprintf("%s-%s-%s", notaID, status, safeFilename)
	dst := filepath.Join(s.dataDir, outputName)

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo de destino: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	return outputName, nil
}

// ListByNotaID lista todos os arquivos de notas fiscais no diretório de dados
// cujo nome possui o prefixo "{notaID}-". Os resultados são ordenados por nome
// em ordem decrescente. Cada arquivo listado contém o nome, o status extraído
// do nome do arquivo e a URL relativa para download. Se o diretório não existir,
// retorna uma lista vazia sem erro.
func (s *NotaService) ListByNotaID(notaID string) ([]models.ListedNota, error) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.ListedNota{}, nil
		}
		return nil, fmt.Errorf("falha ao ler diretório: %w", err)
	}

	prefix := notaID + "-"
	notas := make([]models.ListedNota, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		notas = append(notas, models.ListedNota{
			Name:   name,
			Status: extractStatusFromFilename(name, notaID),
			URL:    fmt.Sprintf("/retrieve/%s", name),
		})
	}

	sort.Slice(notas, func(i, j int) bool {
		return notas[i].Name > notas[j].Name
	})

	return notas, nil
}

// GetFilePath retorna o caminho completo do arquivo de nota fiscal no diretório de dados.
func (s *NotaService) GetFilePath(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

// extractStatusFromFilename extrai o status do nome do arquivo de nota fiscal.
// O formato esperado é "{notaID}-{status}-{resto...}". Retorna o status como string
// ou string vazia se o formato não corresponder ao esperado.
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
