# Telegram Connection Status Checker Bot

Este é um bot Telegram desenvolvido em Go, projetado para verificar o status de conexão e inadimplência de usuários. O bot utiliza o pacote `tgbotapi` para se comunicar com a API do Telegram, oferecendo funcionalidades como comandos básicos e menus interativos.

---

## 📋 Funcionalidades

1. **Comando `/start`:** 
   - Exibe uma mensagem de boas-vindas.
   - Apresenta um menu interativo com opções.

2. **Consulta de inadimplentes:**
   - Verifica e retorna o status de inadimplência.
   - Detalha os dados obtidos, incluindo:
     - Total de inadimplentes.
     - Quantidade de conexões bloqueadas e não bloqueadas.
     - Reduções aplicadas.

3. **Segurança:**
   - Apenas usuários autorizados podem interagir com o bot.
   - Controle para evitar execução simultânea de comandos por um mesmo usuário.

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** Go (Golang)
- **Biblioteca:** [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api)

---

## 🚀 Como Executar

### Pré-requisitos

1. **Instalar o Go:**
   - [Guia de instalação do Go](https://golang.org/doc/install)

2. **Obter o Token do Bot:**
   - Crie um bot no Telegram com o [BotFather](https://core.telegram.org/bots#botfather) e copie o token fornecido.

3. **Definir IDs Autorizados:**
   - Inclua os IDs dos usuários autorizados no mapa `authorizedUserIDs` no código.

---

### Passo a Passo

1. Clone o repositório:
   ```bash
   git clone https://github.com/seu-usuario/telegram-connection-status-checker.git
   cd telegram-connection-status-checker
