# stage 1 (build)
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o central-config-service .

# stage 2 (final)
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/central-config-service .

EXPOSE 8000

CMD ["./central-config-service"]