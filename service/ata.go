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

// AtaService gerencia o armazenamento de arquivos de atas de registro de preço
// e fornece integração com a API do Portal Nacional de Contratações Públicas (PNCP).
type AtaService struct {
	dataDir string
	client  http.Client
	urlPncp string
}

// NewAtaService cria um novo AtaService com o diretório de dados, a URL base da API
// do PNCP e o cliente HTTP para realizar as requisições.
func NewAtaService(dataDir, urlPncp string, client http.Client) *AtaService {
	return &AtaService{
		dataDir: dataDir,
		client:  client,
		urlPncp: urlPncp,
	}
}
// EnsureAtaDataDir garante que o diretório de dados de atas exista,
// criando-o com permissões 0755 se necessário.
func (s *AtaService) EnsureAtaDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

// SaveFile salva um arquivo de ata no diretório de dados com o nome composto
// no formato "{ataId}-{nomeOriginal}". Retorna o nome do arquivo gerado
// ou um erro em caso de falha na abertura, criação ou cópia do arquivo.
func (s *AtaService) SaveFile(fileHeader *multipart.FileHeader, ataId string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer src.Close()

	safeFilename := filepath.Base(fileHeader.Filename)
	outputFilename := fmt.Sprintf("%s-%s", ataId, safeFilename)
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
// GetAta retorna o caminho completo do arquivo de ata no diretório de dados.
func (s *AtaService) GetAta(filename string) string {
	return filepath.Join(s.dataDir, filename)
}
// DeleteAta remove um arquivo de ata do diretório de dados.
// Retorna erro se a remoção falhar.
func (s *AtaService) DeleteAta(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("falha ao remover arquivo: %w", err)
	}
	return nil
}

// GetAtaInfoPncp consulta a API do PNCP para obter informações de uma ata de registro
// de preço específica, identificada pelo CNPJ do órgão, ano, sequencial da compra
// e sequencial da ata. Retorna os dados da ata ou erro em caso de falha na requisição,
// parse da resposta ou ata não encontrada (404).
func (s *AtaService) GetAtaInfoPncp(cnpj string, year string, sequencialCompra string, sequencialAta string) (*models.AtaPncp, error) {

	urlReq := fmt.Sprintf("%s/v1/orgaos/%s/compras/%s/%s/atas/%s",
		s.urlPncp, cnpj, year, sequencialCompra, sequencialAta)

	resp, err := s.client.Get(urlReq)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("ata não encontrada no PNCP (404)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do PNCP retornou status inesperado: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao processar a resposta: %w", err)
	}

	data := &models.AtaPncp{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, fmt.Errorf("falha ao recuperar ata do PNCP: %w", err)
	}

	return data, nil
}
