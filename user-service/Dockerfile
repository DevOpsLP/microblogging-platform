# Base de Go para compilar la aplicación
FROM golang:1.23 as builder

WORKDIR /app

# Instala dependencias
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila la aplicación en un ejecutable binario
RUN go build -o user-service cmd/main.go

# Etapa final con ubuntu:latest
FROM ubuntu:latest

WORKDIR /

COPY --from=builder /app/user-service /user-service

# Expone el puerto que usa user-service
EXPOSE 8080

CMD ["/user-service"]
