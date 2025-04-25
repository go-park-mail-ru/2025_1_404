FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./main ./app && chmod +x ./main

EXPOSE 8001

CMD ["./main"]
