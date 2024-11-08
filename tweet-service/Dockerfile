# Base de Go para compilar la aplicación
FROM golang:1.23 as builder

WORKDIR /app

# Instala dependencias
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila la aplicación en un ejecutable binario
RUN go build -o tweet-service cmd/main.go

# Etapa final con ubuntu:latest
FROM ubuntu:latest

WORKDIR /

COPY --from=builder /app/tweet-service /tweet-service

# Expone el puerto que usa tweet-service
EXPOSE 8081

CMD ["/tweet-service"]
