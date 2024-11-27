package junior

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang-telegram-connection-status-checker/services/graphql"
)

var (
	apiUrl    = os.Getenv("JUNIOR_URL")
	authToken = os.Getenv("JUNIOR_AUTH_TOKEN")
)

func makeRequest(username, ipComunicacao string, logger *log.Logger) (ApiResponse, error) {
	client := &http.Client{}
	data := map[string]string{"username": username, "bng_ip": ipComunicacao}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Printf("Erro ao criar JSON da requisição: %v", err)
		return ApiResponse{}, fmt.Errorf("erro ao criar JSON da requisição: %v", err)
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(string(jsonData)))
	if err != nil {
		logger.Printf("Erro ao criar requisição: %v", err)
		return ApiResponse{}, fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("Erro ao fazer a requisição: %v", err)
		return ApiResponse{}, fmt.Errorf("erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Log da resposta da API
	logger.Printf("Requisição para %s com IP %s retornou status: %s", username, ipComunicacao, resp.Status)

	if resp.StatusCode != http.StatusOK {
		return ApiResponse{}, fmt.Errorf("erro HTTP: %s", resp.Status)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		logger.Printf("Erro ao decodificar resposta JSON: %v", err)
		return ApiResponse{}, fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	// Log da resposta decodificada
	logger.Printf("Resposta da API: Status = %s, Plano = %s", apiResp.Status, apiResp.Plano)

	return apiResp, nil
}

func SaveToCSV(users []graphql.Inadimplente, logger *log.Logger) error {
	file, err := os.Create("inadimplentes.csv")
	if err != nil {
		logger.Printf("Erro ao criar arquivo CSV: %v", err)
		return fmt.Errorf("erro ao criar arquivo CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escreve o cabeçalho no arquivo CSV
	headers := []string{"Username", "IP Comunicação", "Status", "Plano", "Nome/Razão Social", "Nome Revenda"}
	if err := writer.Write(headers); err != nil {
		logger.Printf("Erro ao escrever cabeçalho no arquivo CSV: %v", err)
		return fmt.Errorf("erro ao escrever cabeçalho no arquivo CSV: %v", err)
	}

	// Cria um conjunto (map) para armazenar usernames processados
	processedUsernames := make(map[string]bool)

	// Escreve os dados dos inadimplentes no arquivo CSV
	for i, user := range users {
		// Verifica se o username já foi processado
		if processedUsernames[user.Username] {
			logger.Printf("Username %s já processado, pulando requisição", user.Username)
			continue
		}

		status := "Desconhecido"
		plano := "Desconhecido"

		response, err := makeRequest(user.Username, user.IpComunicacao, logger)
		if err == nil {
			status = response.Status
			plano = response.Plano
		}

		// Marca o username como processado
		processedUsernames[user.Username] = true

		logger.Printf("Escrevendo linha %d com %s, %s", i+2, user.Username, user.IpComunicacao)

		record := []string{user.Username, user.IpComunicacao, status, plano, user.NomeRazaoSocial, user.NomeRevenda}
		if err := writer.Write(record); err != nil {
			logger.Printf("Erro ao escrever linha no arquivo CSV: %v", err)
			return fmt.Errorf("erro ao escrever linha no arquivo CSV: %v", err)
		}
	}

	logger.Println("Arquivo CSV gerado com sucesso!")
	return nil
}
