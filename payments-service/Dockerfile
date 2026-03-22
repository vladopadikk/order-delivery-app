FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment-service ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-service .

CMD ["./payment-service"]