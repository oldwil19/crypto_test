FROM --platform=linux/amd64 golang:1.23.3-alpine3.20 AS builder
RUN apk add --no-cache build-base gcc g++ libc6-compat linux-headers

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copiar el código fuente
COPY . .

# Configuración de compilación
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o server ./cmd/server/main.go

FROM --platform=linux/amd64 alpine:3.20

# Instalar dependencias necesarias quizas no haga falta pero luego ve
RUN apk add --no-cache postgresql-client curl

WORKDIR /app

COPY --from=builder /app/server /app/server

COPY ./migrations /app/migrations
COPY ./configs /app/configs
COPY ./.env /app/.env
COPY ./docs /app/docs
#usare ngixn
EXPOSE 8080 

# Ejecutar el servidor
CMD ["./server"]
