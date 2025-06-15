# Admin Frontend - Support Chat System

This is the admin/operator interface for managing dialogs in the support chat system.

## Features

- Login as Admin or Operator
- View queued and active dialogs
- Real-time chat with clients (via WebSocket)
- Finish/close dialogs
- Operator management (Admins only)
- View statistics (Admins only)

---

## Getting Started

### Prerequisites

- Node.js v18+

### Installation

```bash
cd admin-frontend
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
admin-frontend/
├── src/
│   ├── components/       # UI components like DialogPanel, Sidebar
│   ├── views/            # Pages (Chat, Operators, History)
│   ├── services/         # API and WebSocket logic
│   └── layouts/          # Layouts (AdminLayout)
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

- Admin and operator JWTs are stored in `localStorage`
- WebSocket messages include system messages (e.g. dialog taken/closed)

---

## Authors

- DIMA AVTENEV (https://github.com/HIPPER52)

---

## License

MIT
