FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN go build -o documentstore ./cmd/server

FROM alpine:3.19 AS deploy

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./documentstore"]