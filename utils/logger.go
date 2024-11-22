// utils/logger.go
package utils

import (
	"log"
	"os"
)

// ConfigureLogger configura um logger personalizado
func ConfigureLogger() *log.Logger {
	// Abre o arquivo de log ou cria um novo se não existir
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir ou criar arquivo de log:", err)
	}

	// Configura o logger para escrever no arquivo e também no console
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
