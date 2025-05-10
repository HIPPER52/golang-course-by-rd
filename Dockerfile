FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o documentstore ./cmd/server

FROM alpine:3.19 AS deploy

WORKDIR /app

COPY --from=builder /app/documentstore .

EXPOSE 8080

CMD ["./documentstore"]