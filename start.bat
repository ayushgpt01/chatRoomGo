@echo off
echo Starting App...

:: Start the Go Backend in a new window
start "Backend" cmd /c "go run ./cmd/server/"

:: Change to frontend directory and start Vite in a new window
start "Frontend" cmd /c "cd frontend && npm run dev"

echo Dev servers are running.
pause