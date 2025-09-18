FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

# Копируем configdb.yaml как config.yaml (чтобы приложение читало config.yaml)
COPY config/configdb.yaml /app/config/config.yaml

COPY migrations /app/migrations

EXPOSE 8080

CMD ["./app"]
