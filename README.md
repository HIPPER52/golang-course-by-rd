# Golang Support Chat System

This is a support chat system built in Go, Vue.js, and MongoDB. It includes a backend service, an admin interface, and a client-facing interface.

## Features

- Real-time chat using WebSocket
- Dialog queue management (Queued → Active → Archived)
- Role-based access control (Admin and Operator)
- Admin dashboard to manage users and view statistics
- Message persistence via MongoDB and RabbitMQ
- Dialog history and reporting
- Basic statistics on operator performance

## Services

### 1. Backend

- Built with Go (Fiber v2)
- MongoDB for persistence
- RabbitMQ for message queueing
- REST API for admin and client services
- WebSocket gateway for real-time communication

### 2. Admin Frontend

- Built with Vue 3
- Role-based dashboard for admins and operators
- Operators can view, pick, and respond to dialogs
- Admins can manage users and see dialog stats

### 3. Client Frontend

- Simple chatbox interface for clients
- Connects to backend via WebSocket
- Messages saved and visible to operators

## Prerequisites

Make sure you have the following installed:

- Docker
- Docker Compose
- Make
- Go (for development)
- Node.js and npm (for frontend development)

## Setup Instructions

1. Clone the repository.
2. Run the following commands to set up the environment:

```bash
make setup-envs
make up
