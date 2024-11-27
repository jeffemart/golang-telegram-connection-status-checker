# Etapa 1: Construção da imagem
FROM golang:1.23-alpine AS builder

# Defina o diretório de trabalho
WORKDIR /app

# Copie o go.mod e go.sum para o contêiner e baixe as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código-fonte para o contêiner
COPY . .

# Compile o aplicativo Go
RUN go build -o myapp .

# Etapa 2: Imagem final para rodar o aplicativo
FROM alpine:latest

# Instale as dependências necessárias (caso precise de alguma, por exemplo, libc)
RUN apk --no-cache add ca-certificates

# Copie o binário compilado da imagem builder
COPY --from=builder /app/myapp /usr/local/bin/myapp

# Defina o comando de execução do contêiner
CMD ["myapp"]
