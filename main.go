package main

import (
	"fmt"
	"log"
	"os"

	"golang-telegram-connection-status-checker/services/graphql"
	"golang-telegram-connection-status-checker/services/junior"
	"golang-telegram-connection-status-checker/utils"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// IDs autorizados para interagir com o bot
var authorizedUserIDs = map[int64]bool{
	1441826228: true,
	987654321:  true,
}

// Armazena o estado do comando para cada usuário
var commandStatus = make(map[int64]bool)

func main() {
	// Token do bot (substitua com seu token)
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	// Cria o bot com o token
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Erro ao criar o bot: %v", err)
	}

	// Ativa o modo de depuração
	bot.Debug = false
	log.Printf("Bot autorizado como %s", bot.Self.UserName)

	// Configuração de atualização
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Obtém as atualizações do bot
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Erro ao obter atualizações: %v", err)
	}

	// Loop principal do bot
	for update := range updates {
		if update.Message != nil {
			// Verifica se o ID do usuário está autorizado
			if !isAuthorizedUser(update.Message.From.ID) {
				continue
			}

			// Trata o comando "/start"
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá! Bem-vindo ao bot. Escolha uma opção abaixo:")
				bot.Send(msg)
				showInlineMenu(bot, update.Message)
			}
		}

		// Verifica se há um CallbackQuery
		if update.CallbackQuery != nil {
			// Processa o CallbackQuery
			switch update.CallbackQuery.Data {
			case "/start":
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Bot iniciado!"))
			case "/inadimplentes":
				handleInadimplentes(bot, update.CallbackQuery.Message)
			case "/relatorio":
				handleRelatorio(bot, update.CallbackQuery.Message)
			}

			// Envia uma confirmação do callback
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Comando recebido!")
			if _, err := bot.AnswerCallbackQuery(callback); err != nil {
				log.Printf("Erro ao responder ao callback: %v", err)
			}
		}
	}
}

// Verifica se o usuário está autorizado
func isAuthorizedUser(userID int) bool {
	_, authorized := authorizedUserIDs[int64(userID)]
	return authorized
}

// Exibe o menu com botões inline
func showInlineMenu(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	startButton := tgbotapi.NewInlineKeyboardButtonData("Iniciar o bot", "/start")
	inadimplentesButton := tgbotapi.NewInlineKeyboardButtonData("Verificar inadimplentes", "/inadimplentes")
	relatorioButton := tgbotapi.NewInlineKeyboardButtonData("Gerar relatório", "/relatorio") // Novo botão

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(startButton, inadimplentesButton, relatorioButton), // Inclui o botão de relatorio
	)

	menuMessage := tgbotapi.NewMessage(msg.Chat.ID, "Escolha uma das opções:")
	menuMessage.ReplyMarkup = inlineKeyboard

	if _, err := bot.Send(menuMessage); err != nil {
		log.Printf("Erro ao enviar o menu: %v", err)
	}
}

// Função para tratar a consulta de inadimplentes
func handleInadimplentes(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Println("Comando /inadimplentes recebido")

	userID := int64(msg.From.ID)
	if commandStatus[userID] {
		log.Println("Comando já executado, aguardando...")
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Você já executou este comando. Por favor, aguarde."))
		return
	}

	commandStatus[userID] = true
	defer func() {
		commandStatus[userID] = false
		log.Println("Comando liberado para novo uso.")
	}()

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Consultando inadimplentes, por favor aguarde..."))

	// Consultar inadimplentes
	query := `
	query MyQuery {
		mk01 {
			inadimplentes_45dias {
				codcontrato
				conexao_bloqueada
				esta_reduzida
				ip_comunicacao
				nome_razaosocial
				nome_revenda
				username
			}
		}
	}
	`

	inadimplentes, err := graphql.FetchInadimplentes(query)
	if err != nil {
		log.Printf("Erro ao buscar inadimplentes: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Erro ao consultar inadimplentes."))
		return
	}

	// Verificar se inadimplentes foi retornado corretamente
	if len(inadimplentes) == 0 {
		log.Println("Nenhum inadimplente encontrado.")
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Nenhum inadimplente encontrado."))
		return
	}

	// Contadores
	total := len(inadimplentes)
	bloqueados := 0
	naoBloqueados := 0
	reduzidos := 0
	naoReduzidos := 0

	// Contar os estados
	for _, inad := range inadimplentes {
		if inad.ConexaoBloqueada == "" {
			inad.ConexaoBloqueada = "Desconhecido"
		}
		if inad.EstaReduzida == "" {
			inad.EstaReduzida = "Desconhecido"
		}

		switch inad.ConexaoBloqueada {
		case "S":
			bloqueados++
		case "N":
			naoBloqueados++
		default:
			log.Printf("Valor inesperado em ConexaoBloqueada: %s", inad.ConexaoBloqueada)
		}

		switch inad.EstaReduzida {
		case "S":
			reduzidos++
		case "N":
			naoReduzidos++
		default:
			log.Printf("Valor inesperado em EstaReduzida: %s", inad.EstaReduzida)
		}
	}

	// Formatar a resposta
	response := fmt.Sprintf(
		"Resumo dos inadimplentes de 45 dias:\n\n"+
			"Total de inadimplentes: %d\n"+
			"Conexão Bloqueada (S): %d\n"+
			"Conexão Não Bloqueada (N): %d\n"+
			"Reduzidos (S): %d\n"+
			"Não Reduzidos (N): %d\n",
		total, bloqueados, naoBloqueados, reduzidos, naoReduzidos,
	)

	// Enviar a resposta
	msgFinal := tgbotapi.NewMessage(msg.Chat.ID, response)
	_, err = bot.Send(msgFinal)
	if err != nil {
		log.Printf("Erro ao enviar a mensagem final: %v", err)
	}
}

// Função para tratar o comando de gerar o relatório
func handleRelatorio(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Println("Comando /relatorio recebido")
	logger := utils.ConfigureLogger()

	userID := int64(msg.From.ID)
	if commandStatus[userID] {
		log.Println("Comando já executado, aguardando...")
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Você já executou este comando. Por favor, aguarde."))
		return
	}

	commandStatus[userID] = true
	defer func() {
		commandStatus[userID] = false
		log.Println("Comando liberado para novo uso.")
	}()

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Gerando relatório, por favor aguarde..."))

	// Consultar inadimplentes através do GraphQL para 30 dias
	query30 := `
	query MyQuery {
		mk01 {
			inadimplentes_30dias {
				codcontrato
				conexao_bloqueada
				esta_reduzida
				ip_comunicacao
				nome_razaosocial
				nome_revenda
				username
			}
		}
	}
	`
	inadimplentes30, err := graphql.FetchInadimplentes(query30)
	if err != nil {
		log.Printf("Erro ao consultar inadimplentes 30 dias: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Erro ao consultar inadimplentes de 30 dias."))
		return
	}

	// Consultar inadimplentes através do GraphQL para 45 dias
	query45 := `
	query MyQuery {
		mk01 {
			inadimplentes_45dias {
				codcontrato
				conexao_bloqueada
				esta_reduzida
				ip_comunicacao
				nome_razaosocial
				nome_revenda
				username
			}
		}
	}
	`
	inadimplentes45, err := graphql.FetchInadimplentes(query45)
	if err != nil {
		log.Printf("Erro ao consultar inadimplentes 45 dias: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Erro ao consultar inadimplentes de 45 dias."))
		return
	}

	// Combine as listas de inadimplentes
	inadimplentes := append(inadimplentes30, inadimplentes45...)

	// Salvar os dados em um arquivo CSV
	if err := junior.SaveToCSV(inadimplentes, logger); err != nil {
		log.Printf("Erro ao salvar o CSV: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Erro ao gerar o relatório CSV."))
		return
	}

	// Enviar o arquivo CSV gerado
	file := tgbotapi.NewDocumentUpload(msg.Chat.ID, "inadimplentes.csv")
	_, err = bot.Send(file)
	if err != nil {
		log.Printf("Erro ao enviar o arquivo CSV: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Erro ao enviar o relatório CSV."))
		return
	}

	// Enviar confirmação de sucesso
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Relatório gerado com sucesso!"))
}
