.PHONY: build up up-dev run down clean deps

APP_NAME=documentstore
DOCKER_COMPOSE=docker-compose
PORT=8080

build:
	@echo "Building Go binary..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd/server

up:
	@echo "Starting full stack in Docker..."
	$(DOCKER_COMPOSE) up --build

up-dev:
	@echo "Starting dependencies only (MongoDB, RabbitMQ)..."
	$(DOCKER_COMPOSE) up mongo rabbitmq -d

run:
	@echo "Running server locally..."
	go run ./cmd/server/main.go

down:
	@echo "Stopping and removing containers..."
	$(DOCKER_COMPOSE) down

clean:
	@echo "Cleaning binary..."
	rm -rf bin/$(APP_NAME)

deps:
	go mod tidy
	go mod download