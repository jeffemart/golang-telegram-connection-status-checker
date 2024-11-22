// services/graphql/graphql.go
package graphql

import (
	"os"
	"bytes"
	"encoding/json"
	"fmt"
	"golang-telegram-connection-status-checker/utils"
	"net/http"
	"time"
)

// Função para fazer a requisição GraphQL
func FetchInadimplentes(query string) ([]Inadimplente, error) {
	// Obtém o logger configurado do pacote utils
	logger := utils.ConfigureLogger()

	// URL do endpoint GraphQL
	graphQLURL := os.Getenv("GRAPHQL_URL")

	if graphQLURL == "" {
		logger.Println("GRAPHQL_URL is not set")
		return nil, fmt.Errorf("GRAPHQL_URL is not set")
	}

	// Cabeçalhos para autenticação
	headers := map[string]string{
		"Content-Type":         "application/json",
		"x-hasura-admin-secret": os.Getenv("HASURA_SECRET"),
	}

	// Corpo da requisição
	requestBody := map[string]string{
		"query": query,
	}

	// Converte o corpo da requisição para JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		logger.Println("Erro ao converter corpo da requisição para JSON:", err)
		return nil, fmt.Errorf("erro ao converter corpo da requisição para JSON: %v", err)
	}

	// Faz a requisição HTTP
	req, err := http.NewRequest("POST", graphQLURL, bytes.NewBuffer(body))
	if err != nil {
		logger.Println("Erro ao criar requisição HTTP:", err)
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %v", err)
	}

	// Adiciona os cabeçalhos
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Faz a requisição HTTP
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println("Erro ao fazer a requisição:", err)
		return nil, fmt.Errorf("erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Decodifica a resposta JSON
	var data ResponseData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Println("Erro ao decodificar resposta JSON:", err)
		return nil, fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	// Retorna os inadimplentes de acordo com a consulta
	return data.Data.Mk01.Inadimplentes30Dias, nil
}