# Chat Room Web Server with WebSocket

This project implements a simple chat room web server using WebSocket for real-time communication between clients.

## Features

- **WebSocket Communication:** Allows real-time messaging between clients.
- **Message Handling:** Supports sending and receiving messages with basic error handling.
- **Static File Serving:** Serves static HTML and JavaScript files for client-side rendering.
- **Custom Chat Components:** Includes custom web components for rendering chat messages.

## Requirements

- Go (Golang)
- WebSocket library (Gorilla WebSocket for Go, native WebSocket API for JavaScript)

## Setup

### Server (Go)

1. Clone the repository:

   ```bash
   git clone https://github.com/ayushgpt01/chatRoomGo.git
   cd chat-room
   ```

2. Build and run the server:

   ```bash
   go run main.go
   ```

The server will start at http://localhost:8080.

<!-- ### Client (JavaScript)

1. Ensure Node.js is installed.

2. Install dependencies:

   ```bash
   npm install
    ```

3. Start the development server:

    ```bash
    npm start
    ```

This will serve the client at http://localhost:3000. -->

## Usage

1. Open `http://localhost:8000` in your web browser.
2. Enter a message in the input field and click "Send Message" to send it to other connected clients.
3. Messages from other clients will appear in the chat window.

## Project Structure

```
main.go # Server file
go.mod
go.sum
static
├── index.html   # HTML template for the chat room interface
├── index.js     # JavaScript for handling UI and WebSocket connections
└── utils
    ├── helpers.js          # JavaScript helper functions
    ├── websocket.js        # WebSocket initialization and event handling
    └── components
        ├── CustomChatElement.js  # Custom web component for chat messages
        ├── LeftMessage.js         # Left side chat message component
        └── RightMessage.js        # Right side chat message component
```
