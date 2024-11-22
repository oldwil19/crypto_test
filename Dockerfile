# Etapa 1: Builder
FROM --platform=linux/amd64 golang:1.23.3-alpine3.20 AS builder

# Instalar herramientas necesarias
RUN apk add --no-cache build-base gcc g++ libc6-compat linux-headers git

# Configuración del entorno de trabajo
WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copiar el código fuente
COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o server ./cmd/server/main.go

# Etapa 2: Imagen final
FROM --platform=linux/amd64 alpine:3.20

# Instalar dependencias necesarias
RUN apk add --no-cache postgresql-client curl

# Configuración del entorno de trabajo
WORKDIR /app

# Copiar binario de la aplicación desde la etapa de compilación
COPY --from=builder /app/server /app/server

# Copiar archivos necesarios
COPY ./migrations /app/migrations
COPY ./configs /app/configs
COPY ./.env /app/.env
COPY ./docs /app/docs

# Exponer el puerto de la aplicación
EXPOSE 8080

# Comando para ejecutar el servidor
CMD ["./server"]
