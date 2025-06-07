# Backend - Support Chat System

This is the backend service for a real-time support chat application. It handles:

- Client/operator/admin registration and authorization
- Dialog queue management
- WebSocket real-time messaging
- Persistence via MongoDB
- Messaging via RabbitMQ

## Features

- REST API (Fiber v2)
- WebSocket gateway for operators and clients
- JWT authentication
- Message queuing and delivery via RabbitMQ
- Unit tests for services

---

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### Installation

```bash
git clone https://github.com/HIPPER52/golang-course-by-rd
cd golang-course-by-rd/backend
go mod tidy
```

### Run with Docker Compose

```bash
make up       # Runs MongoDB and RabbitMQ
make run-server     # Runs backend server locally
make run-consumer   # Runs message consumer
```

Or run everything inside Docker:

```bash
make up
```

### Build binary

```bash
make build
```

---

## Project Structure

```
backend/
├── cmd/                  # Entry points (server, consumer)
├── internal/
│   ├── api/              # HTTP handlers
│   ├── services/         # Business logic
│   ├── repository/       # DB interactions
│   ├── models/           # DTOs and Mongo models
│   ├── ws/               # WebSocket gateway
│   ├── consumer/         # RabbitMQ consumers
│   ├── producer/         # RabbitMQ producers
│   └── constants/        # Constants and enums
├── docker-compose.yml
└── Makefile
```

---

## Environment Variables

Set the following environment variables (you can use `.env` file or pass manually):

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=chat_support
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
JWT_SECRET=your_secret
```

---

## Testing

Run unit tests:

```bash
make test
```

---

## Authors

- DIMA AVTENEV (https://github.com/HIPPER52)

---

## License

MIT
