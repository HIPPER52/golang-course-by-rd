# Client Frontend - Support Chat System

This is the public-facing client interface for initiating chat with support agents.

## Features

- Enter client info (name, phone)
- Start chat dialog
- Real-time messaging via WebSocket
- Persisted message history (from backend)
- Automatically reconnects

---

## Getting Started

### Prerequisites

- Node.js v18+

### Installation

```bash
cd client-frontend
npm install
```

### Development Mode

```bash
npm run dev
```

### Linting & Formatting

```bash
npm run lint
npm run lint:fix
```

### Build

```bash
npm run build
```

---

## Project Structure

```
client-frontend/
├── src/
│   ├── views/            # Main Chat view
│   ├── components/       # (Optional) extracted UI components
│   ├── services/         # API and socket logic
│   └── socket.js         # WebSocket connection handler
├── public/
├── vite.config.js
└── eslint.config.js
```

---

## Environment Variables

For development:

```bash
VITE_API_BASE=http://localhost:8080
```

Create a `.env` file if needed.

---

## Notes

- Client ID and dialog ID are stored in localStorage
- System messages from operator (e.g., dialog taken/closed) are displayed

---

## Authors

- DIMA AVTENEV (https://github.com/HIPPER52)

---

## License

MIT
