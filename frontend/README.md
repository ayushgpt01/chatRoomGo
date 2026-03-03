# ChatRoom Frontend

A modern, type-safe React web application for real-time chatting. Built with **Vite**, **React 19**, and the **TanStack** ecosystem.

## Tech Stack

- **Framework**: [React 19](https://react.dev/)
- **Routing**: [TanStack Router](https://tanstack.com/router) (File-based, type-safe routing)
- **State Management**: 
  - Server State: [TanStack Query v5](https://tanstack.com/query)
  - Client State: [Zustand](https://github.com/pmndrs/zustand)
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/) & [DaisyUI](https://daisyui.com/)
- **Validation**: [Zod](https://zod.dev/) (Schema-based validation for API responses)
- **Tooling**: [Biome](https://biomejs.dev/) (Fast linting and formatting)

## Architecture

The application follows a modular structure focused on type safety and separation of concerns:

- **`/src/routes`**: Contains the file-based route definitions.
- **`/src/services`**: API abstraction layer using Axios.
- **`/src/stores`**: Global client-side state (Auth, Socket, Toasts).
- **`/src/integrations`**: Configuration for external libraries like Axios and TanStack Query.
- **`/src/types`**: Centralized TypeScript interfaces and Zod schemas.


## Getting Started

### Prerequisites
- Node.js (Latest LTS recommended)
- The Go Backend running on `http://localhost:8080`

### Installation
1. Navigate to the frontend directory:
   ```bash
   cd frontend

```

2. Install dependencies:
```bash
npm install

```


3. Start the development server:
```bash
npm run dev

```


The app will be available at `http://localhost:3000`.

## Available Scripts

* `npm run dev`: Starts the Vite development server.
* `npm run build`: Compiles the application for production.
* `npm run check`: Runs Biome linting and formatting checks.
* `npm run test`: Executes unit tests via Vitest.

## Authentication Flow

The app uses a hybrid approach:

1. **Zustand (`authStore`)**: Manages the local user session and tokens.
2. **Axios Interceptor**: Automatically attaches JWT tokens to outgoing requests.
3. **Zod Validation**: All user data returned from the `/auth` endpoints is validated at runtime to prevent malformed data from crashing the UI.

## Real-time Communication

WebSocket connections are managed through a dedicated `socketStore`. This allows the UI to remain reactive to incoming `message` and `room` events regardless of which component the user is currently viewing.
