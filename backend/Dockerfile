FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o /tmp/bin/server ./cmd/server && chmod +x /tmp/bin/server

FROM alpine:3.19 AS deploy

WORKDIR /app

COPY --from=builder /tmp/bin/server /usr/local/bin

EXPOSE 8080

CMD ["/usr/local/bin/server"]