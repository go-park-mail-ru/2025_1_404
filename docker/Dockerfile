FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/auth ./cmd/auth/main.go && chmod +x /app/auth
RUN go build -o /app/offer ./cmd/offer/main.go && chmod +x /app/offer
RUN go build -o /app/zhk ./cmd/zhk/main.go && chmod +x /app/zhk