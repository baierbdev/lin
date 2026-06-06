package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// AditivoService gerencia o armazenamento de arquivos de aditivos contratuais em disco.
// Os arquivos são salvos no diretório de dados com nomes compostos para identificação.
type AditivoService struct {
	dataDir string
}

// NewAditivoService cria um novo AditivoService que utilizará o diretório de dados especificado.
func NewAditivoService(dataDir string) *AditivoService {
	return &AditivoService{
		dataDir: dataDir,
	}
}

// EnsureAditivoDataDir garante que o diretório de dados de aditivos exista,
// criando-o com permissões 0755 se necessário.
func (s *AditivoService) EnsureAditivoDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

// SaveFile salva um arquivo de aditivo no diretório de dados com o nome composto
// no formato "{contratoId}-{tipo}-{date}-{nomeOriginal}". Retorna o nome do arquivo
// gerado ou um erro em caso de falha na abertura, criação ou cópia do arquivo.
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

// GetAditivo retorna o caminho completo do arquivo de aditivo no diretório de dados.
func (s *AditivoService) GetAditivo(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

// DeleteAditivo remove um arquivo de aditivo do diretório de dados.
// Retorna erro se a remoção falhar.
func (s *AditivoService) DeleteAditivo(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("falha ao remover arquivo: %w", err)
	}
	return nil
}
