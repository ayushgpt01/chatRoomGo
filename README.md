# ChatRoomGo

I built this project to explore how to handle real-time, concurrent connections in **Go** and how to manage that state effectively on the **Frontend** using the modern TanStack ecosystem.

It’s a multi-room chat app where everything—from user accounts to message history—is persisted in a local SQLite database.

## 💡 What’s interesting about this project?

### **1. The Backend (Go)**

* **Real-time with WebSockets:** I built a `Hub` system for websockets implementation. It uses Go routines and channels to broadcast messages to the right rooms without blocking the main execution.
* **CGO-Free SQLite:** I used a pure Go SQLite driver (`modernc.org/sqlite`). This means the project is super easy to compile and run on any machine without needing a C compiler installed.
* **Clean Architecture:** I separated the code into `internal/auth`, `internal/room`, and `internal/ws`. This keeps the business logic away from the database code, making it much easier to maintain or swap out the database later.

### **2. The Frontend (React 19)**

* **Type-Safe Routing:** I used **TanStack Router**. If you try to navigate to a route that doesn't exist or pass the wrong params, the TypeScript compiler will catch it before you even run the code.
* **Smart Data Fetching:** **TanStack Query** handles the API states (loading, error, caching), so the UI always feels snappy and stays in sync with the backend.
* **Validation:** I use **Zod** to validate API responses. If the backend sends something unexpected, the app catches it at the "border" instead of crashing deep inside a component.

## Tech Stack

* **Backend:** Go, Gorilla WebSockets (for real-time), JWT (for auth).
* **Database:** SQLite (Single file, no setup required).
* **Frontend:** React 19, TypeScript, TanStack Router & Query, Tailwind CSS 4.0.
* **State Management:** Zustand (for auth and socket state).

## How it's organized

```text
├── cmd/server          # Entry point where the server starts
├── internal/           # The "brains" of the app
│   ├── auth            # Login/Signup logic
│   ├── ws              # WebSocket hub and client "pumps"
│   └── models          # Database schemas and shared types
├── frontend/           # Vite + React app
│   ├── src/routes      # File-based routing
│   └── src/services    # API calls using Axios

```

## How to run it locally

### **1. Backend**

From the root directory:

```bash
go mod download
go run cmd/server/main.go

```

The server starts on `http://localhost:8080`.

### **2. Frontend**

In a new terminal:

```bash
cd frontend
npm install
npm run dev

```

Open `http://localhost:3000` and you're good to go!
