package models

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

type Orgao struct {
	Cnpj        string `json:"cnpj"`
	RazaoSocial string `json:"razaoSocial"`
	EsferaId    string `json:"esferaId"`
	PoderId     string `json:"poderId"`
}

type Unidade struct {
	CodigoUnidade string `json:"codigoUnidade"`
	NomeUnidade   string `json:"nomeUnidade"`
	MunicipioNome string `json:"municipioNome"`
	CodigoIbge    string `json:"codigoIbge"`
	UfSigla       string `json:"ufSigla"`
	UfNome        string `json:"ufNome"`
}

type OrgaoEntidade Orgao
type OrgaoSubRogado Orgao
type UnidadeSubRogada Unidade
type UnidadeOrgao Unidade
