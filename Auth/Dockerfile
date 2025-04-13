FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app


COPY --from=builder /app/.env ./.env
COPY --from=builder /app/config/config.yaml ./config.yaml
COPY --from=builder /app/auth ./auth

USER nonroot:nonroot

CMD ["./auth", "-config", "./config.yaml"]