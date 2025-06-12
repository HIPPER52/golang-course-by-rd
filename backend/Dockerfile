FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY go.mod go.sum ./

RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o /tmp/bin/server ./cmd/server && \
    go build -o /tmp/bin/consumer ./cmd/consumer && \
    chmod +x /tmp/bin/server /tmp/bin/consumer

FROM alpine:3.19 AS deploy

WORKDIR /app
RUN apk add --no-cache netcat-openbsd
COPY --from=builder /tmp/bin/server /usr/local/bin/server

COPY wait-for.sh /usr/local/bin/wait-for.sh

RUN chmod +x /usr/local/bin/wait-for.sh

EXPOSE 8080
ENTRYPOINT ["sh", "-c", "/usr/local/bin/wait-for.sh rabbitmq:5672 && exec /usr/local/bin/server"]

FROM alpine:3.19 AS consumer

WORKDIR /app
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /tmp/bin/consumer /usr/local/bin/consumer
COPY wait-for.sh /usr/local/bin/wait-for.sh

RUN chmod +x /usr/local/bin/wait-for.sh

ENTRYPOINT ["sh", "-c", "/usr/local/bin/wait-for.sh rabbitmq:5672 && exec /usr/local/bin/consumer"]
