FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/config/config.yaml ./config.yaml
COPY --from=builder /app/bin/app ./app

CMD ["./app", "-config", "./config.yaml"]