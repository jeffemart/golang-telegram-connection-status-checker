// services/graphql/types.go
package graphql

// Struct para representar cada inadimplente
type Inadimplente struct {
	CodContrato      int    `json:"codcontrato"`
	ConexaoBloqueada string `json:"conexao_bloqueada"` // Mudado para string para corresponder ao valor "N"
	EstaReduzida     string `json:"esta_reduzida"`     // Mudado para string para corresponder ao valor "N"
	IpComunicacao    string `json:"ip_comunicacao"`
	NomeRazaoSocial  string `json:"nome_razaosocial"`
	NomeRevenda      string `json:"nome_revenda"`
	Username         string `json:"username"`
}

// Struct que representa a resposta completa da requisição GraphQL
type ResponseData struct {
	Data struct {
		Mk01 struct {
			Inadimplentes30Dias []Inadimplente `json:"inadimplentes_45dias"`
		} `json:"mk01"`
	} `json:"data"`
}
