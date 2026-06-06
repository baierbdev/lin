package models

// ContratoPncp representa um contrato obtido da API do Portal Nacional de
// Contratações Públicas (PNCP), com todos os seus metadados e valores financeiros.
type ContratoPncp struct {
	AnoContrato                 int32            `json:"anoContrato"`
	TipoContrato                TipoContrato     `json:"tipoContrato"`
	NumeroContratoEmpenho       string           `json:"numeroContratoEmpenho"`
	OrgaoEntidade               OrgaoEntidade    `json:"orgaoEntidade"`
	DataAssinatura              string           `json:"dataAssinatura"`
	DataVigenciaInicio          string           `json:"dataVigenciaInicio"`
	DataVigenciaFim             string           `json:"dataVigenciaFim"`
	FrutoAdesao                 string           `json:"frutoAdesao"`
	DataAtualizacao             string           `json:"dataAtualizacao"`
	Nifornecedor                string           `json:"niFornecedor"`
	TipoPessoa                  string           `json:"tipoPessoa"`
	NomeRazaoSocialFornecedor   string           `json:"nomeRazaoSocialFornecedor"`
	NiFornecedorSubContrato     string           `json:"niFornecedorSubContratado"`
	DataPublicacaoPncp          string           `json:"dataPublicacaoPncp"`
	InformacaoComplementar      string           `json:"informacaoComplementar"`
	OrgaoSubRogado              OrgaoSubRogado   `json:"orgaoSubRogado"`
	UnidadeOrgao                UnidadeOrgao     `json:"unidadeOrgao"`
	UnidadeSubRogada            UnidadeSubRogada `json:"unidadeSubRogada"`
	SequencialContrato          int32            `json:"sequencialContrato"`
	Processo                    string           `json:"processo"`
	TipoPessoaSubContratada     string           `json:"tipoPessoaSubContratada"`
	NumeroRetificacao           int32            `json:"numeroRetificacao"`
	NumeroControlePncp          string           `json:"numeroControlePNCP"`
	Receita                     bool             `json:"receita"`
	NumeroParcela               int32            `json:"numeroParcelas"`
	TemRemanejamento            bool             `json:"temRemanejamento"`
	EmendaParlamentar           bool             `json:"emendaParlamentar"`
	NomeFornecedorSubContratado string           `json:"nomeFornecedorSubContratado"`
	ObjetoContrato              string           `json:"objetoContrato"`
	ValorInicial                float64          `json:"valorInicial"`
	ValorParcela                float64          `json:"valorParcela"`
	ValorGlobal                 float64          `json:"valorGlobal"`
	ValorAcumulado              float64          `json:"valorAcumulado"`
	DataAtualizacaoGlobal       string           `json:"dataAtualizacaoGlobal"`
	IndentificacaoCipi          string           `json:"identificadorCipi"`
	UrlCipi                     string           `json:"urlCipi"`
	UsuarioNome                 string           `json:"usuarioNome"`
	CodigoPaisFornecedor        string           `json:"codigoPaisFornecedor"`
	NumeroControlePncpCompras   string           `json:"numeroControlePncpCompra"`
	NumeroControlePncpAta       string           `json:"numeroControlePncpAta"`
}

// AtaPncp representa uma ata de registro de preço obtida da API do PNCP,
// contendo dados da ata, vigência, órgãos envolvidos e informações da compra.
type AtaPncp struct {
	NumeroAta                    string           `json:"numeroAtaRegistroPreco"`
	AnoData                      int32            `json:"anoAta"`
	DataAssinatura               string           `json:"dataAssinatura"`
	DataVigenciaInicio           string           `json:"dataVigenciaInicio"`
	DataVigenciaFim              string           `json:"dataVigenciaFim"`
	DataCancelamento             string           `json:"dataCancelamento"`
	Cancelado                    bool             `json:"cancelado"`
	DataPublicacaoPncp           string           `json:"dataPublicacaoPncp"`
	DataInclusao                 string           `json:"dataInclusao"`
	DataAtualizacao              string           `json:"dataAtualizacao"`
	DataAtualizacaoGlobal        string           `json:"dataAtualizacaoGlobal"`
	SequencialAta                int32            `json:"sequencialAta"`
	NumeroControlePncp           string           `json:"numeroControlePNCP"`
	OrgaoEntidade                OrgaoEntidade    `json:"orgaoEntidade"`
	OrgaoSubRogado               OrgaoSubRogado   `json:"orgaoSubRogado"`
	UnidadeOrgao                 UnidadeOrgao     `json:"unidadeOrgao"`
	UnidadeSubRogada             UnidadeSubRogada `json:"unidadeSubRogada"`
	ModalidadeNome               string           `json:"modalidadeNome"`
	ObjetoCompra                 string           `json:"objetoCompra"`
	InformacaoComplementarCompra string           `json:"informacaoComplementarCompra"`
	UsuarioNome                  string           `json:"usuarioNome"`
	NumeroControlePncpCompra     string           `json:"numeroControlePncpCompra"`
	PossibilidadeAdesao          bool             `json:"possibilidadeAdesao"`
}

// Orgao representa um órgão público com CNPJ, razão social, esfera e poder.
type Orgao struct {
	Cnpj        string `json:"cnpj"`
	RazaoSocial string `json:"razaoSocial"`
	EsferaId    string `json:"esferaId"`
	PoderId     string `json:"poderId"`
}

// Unidade representa uma unidade organizacional com código, nome,
// município (nome e código IBGE) e UF (sigla e nome).
type Unidade struct {
	CodigoUnidade string `json:"codigoUnidade"`
	NomeUnidade   string `json:"nomeUnidade"`
	MunicipioNome string `json:"municipioNome"`
	CodigoIbge    string `json:"codigoIbge"`
	UfSigla       string `json:"ufSigla"`
	UfNome        string `json:"ufNome"`
}

// TipoContrato representa um tipo de contrato cadastrado no PNCP.
type TipoContrato struct {
	Id              int32  `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	DataInclusao    string `json:"dataInclusao"`
	DataAtualizacao string `json:"dataAtualizacao"`
	Status          bool   `json:"statusAtivo"`
}

// CategoriaProcesso representa uma categoria de processo no PNCP.
type CategoriaProcesso struct {
	Id   int32  `json:"id"`
	Nome string `json:"nome"`
}

type OrgaoEntidade Orgao
type OrgaoSubRogado Orgao
type UnidadeSubRogada Unidade
type UnidadeOrgao Unidade
