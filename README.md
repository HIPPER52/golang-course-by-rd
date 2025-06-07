# Golang Support Chat System

This is a support chat system built with Go, Vue.js, MongoDB, and RabbitMQ. The system includes:

* A backend API and WebSocket server written in Go (Fiber v3)
* An admin panel frontend (Vue 3) for managing dialogs and operators
* A client-facing frontend (Vue 3) for starting and continuing chats

## Features

* Real-time messaging with WebSocket
* Dialog queue management: queued → active → archived
* Admin/Operator roles with different permissions
* Chat history browsing
* Operator performance statistics
* System messages (e.g., dialog taken/closed)

## Requirements

* Go (1.21+ recommended)
* Docker & Docker Compose
* Make
* Node.js (18+) and npm (9+) for frontend apps

## Getting Started

### 1. Setup

To prepare the environment, run:

```bash
make setup-envs
make up
```

This will:

* Create necessary `.env` files from examples
* Launch MongoDB, RabbitMQ, and all services via Docker Compose

### 2. Accessing the system

* **Admin frontend:** [http://localhost:5173](http://localhost:5173)
* **Client frontend:** [http://localhost:5174](http://localhost:5174)
* **Backend API (REST/WebSocket):** [http://localhost:8080](http://localhost:8080)

### Admin Credentials

* **Email:** [admin@admin.com](mailto:admin@admin.com)
* **Password:** 12345678

## Development

You can work with services individually outside Docker:

### Backend

```bash
cd backend
go run main.go
```

### Admin Frontend

```bash
cd admin-frontend
npm install
npm run dev
```

### Client Frontend

```bash
cd client-frontend
npm install
npm run dev
```

## Project Structure

```
.
├── backend/             # Go backend (Fiber)
├── admin-frontend/      # Admin interface (Vue 3)
├── client-frontend/     # Client chat UI (Vue 3)
├── docker-compose.yml   # Docker services: MongoDB, RabbitMQ, backend, frontends
├── Makefile             # Development commands
├── README.md            # This file
```

## License

MIT
