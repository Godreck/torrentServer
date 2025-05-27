# Dockerfile

# Стадия сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/torrentServer

# Финальный минимальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

ENV CONFIG_PATH=internal/config/local.yaml

WORKDIR /

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]
