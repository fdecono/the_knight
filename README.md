# The Knight's Tour

A Go implementation of the Knight's Tour problem using Warnsdorff's heuristic algorithm with backtracking.

### Demo

Watch the Knight's Tour solver in action:

<div align="center">
  <img src="the_knight.gif" width="50%" height="50%" alt="Knight's Tour Demo">
</div>

## Overview

The Knight's Tour is a classic chess problem: find a sequence of moves for a knight such that it visits every square on a chessboard exactly once. This implementation solves the problem using:

- **Warnsdorff's Heuristic**: Prioritizes moves that lead to positions with fewer available next moves
- **Backtracking**: Explores all possible paths and backtracks when no solution is found
- **Recursive Algorithm**: Efficiently searches the solution space

### Features

- ✅ Solves the Knight's Tour on an 8×8 chessboard
- ✅ Uses Warnsdorff's heuristic for fast solution finding
- ✅ Interactive Web UI with real-time visualization
- ✅ Execution timing and attempt counting
- ✅ Backtracking support for complete search

## Installation and Usage

### Prerequisites

- Go 1.21 or higher
- Modern web browser with JavaScript enabled

### Installation

```bash
git clone <repository-url>
cd the_knight
go mod tidy
```

### Running the Web Server

Start the web server:

```bash
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o knight-tour ./cmd/server
./knight-tour
```

Then open your browser and visit:
```
http://localhost:8080
```

## Project Structure

```
the_knight/
├── cmd/
│   └── server/
│       └── main.go          # Server entry point
├── internal/
│   ├── solver/
│   │   ├── solver.go        # Core solving algorithm with channels
│   │   └── types.go         # MoveUpdate and SolveResult types
│   └── web/
│       └── server.go        # HTTP server and handlers
├── pkg/
│   └── board/
│       └── board.go         # Board logic (reusable package)
├── web/
│   ├── templates/
│   │   └── index.html      # HTMX frontend
│   └── static/              # Static assets
├── main.go                  # Legacy entry point (now starts web server)
└── go.mod
```

## Architecture

### Concurrency and Channels

The solver uses Go channels for real-time communication between the solving algorithm and the web server:

**Channel-Based Communication:**
- **Buffered Channel (size 1000)**: Prevents solver from blocking when sending updates
- **Non-blocking Sends**: Uses `select` with `default` to avoid deadlocks
- **WaitGroup Coordination**: Ensures goroutine completion before returning results
- **Context Cancellation**: Allows graceful shutdown and request cancellation

**Data Flow:**
```
Solver (Goroutine) → moveChan → Web Server → SSE → Frontend (HTMX)
```

**Concurrency Patterns:**
- Solver runs in a background goroutine
- Moves are sent through a buffered channel to the web server
- Server-Sent Events (SSE) stream moves to the frontend in real-time
- Mutex-protected state for thread-safe access to results
- Channels are drained on failure to prevent memory leaks

**State Management:**
- Moves stored in slice only during solving
- On success: Moves copied to result and retained
- On failure: Moves cleared immediately, channels drained

## License

This project is open source and available for educational purposes.
