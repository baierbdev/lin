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

// ContratoService gerencia o armazenamento de arquivos de contratos em disco
// e fornece integração com a API do Portal Nacional de Contratações Públicas (PNCP).
type ContratoService struct {
	dataDir string
	client  http.Client
	urlPncp string
}

// NewContratoService cria um novo ContratoService com o diretório de dados, a URL base
// da API do PNCP e o cliente HTTP para realizar as requisições.
func NewContratoService(dataDir string, urlPncp string, client http.Client) *ContratoService {
	return &ContratoService{
		dataDir: dataDir,
		client:  client,
		urlPncp: urlPncp,
	}
}

// EnsureContratoDataDir garante que o diretório de dados de contratos exista,
// criando-o com permissões 0755 se necessário.
func (s *ContratoService) EnsureContratoDataDir() error {
	return os.MkdirAll(s.dataDir, 0o755)
}

// SaveFile salva um arquivo de contrato no diretório de dados com o nome composto
// no formato "{contratoId}-{nomeOriginal}". Retorna o nome do arquivo gerado
// ou um erro em caso de falha na abertura, criação ou cópia do arquivo.
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

// GetContrato retorna o caminho completo do arquivo de contrato no diretório de dados.
func (s *ContratoService) GetContrato(filename string) string {
	return filepath.Join(s.dataDir, filename)
}

// Deletecontrato remove um arquivo de contrato do diretório de dados.
// Retorna erro se a remoção falhar.
func (s *ContratoService) Deletecontrato(filename string) error {
	dst := filepath.Join(s.dataDir, filename)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("falha ao remover arquivo: %w", err)
	}
	return nil
}

// GetContratoPncp consulta a API do PNCP para obter informações de um contrato
// específico, identificado pelo CNPJ do órgão, ano e sequencial do contrato.
// Retorna os dados do contrato ou erro em caso de falha na requisição,
// parse da resposta ou contrato não encontrado (404).
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
