# Architecture Overview

## Production-Grade Design Decisions

This implementation follows production-ready patterns with proper separation of concerns, concurrency safety, and scalable architecture.

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
├── main.go                  # Legacy entry point (now redirects to server)
└── go.mod
```

## Key Components

### 1. Channel-Based Communication (`internal/solver/`)

**Design Rationale:**
- **Buffered Channel (size 1000)**: Prevents solver from blocking when sending updates
- **Non-blocking Sends**: Uses `select` with `default` to avoid deadlocks
- **WaitGroup Coordination**: Ensures goroutine completion before returning results
- **Context Cancellation**: Allows graceful shutdown and request cancellation

**MoveUpdate Flow:**
```
Solver → moveChan → Web Server → SSE → Frontend (HTMX)
```

**State Management:**
- Moves stored in slice (`s.moves`) only during solving
- **On Success**: Moves copied to `SolveResult` and retained
- **On Failure**: Moves cleared immediately (memory efficiency)
- Channels drained on failure to prevent stale data

### 2. Web Server (`internal/web/`)

**HTTP Endpoints:**
- `GET /` - Serves HTML with HTMX
- `POST /api/solve` - Starts solving (returns immediately)
- `GET /api/moves/stream` - Server-Sent Events (SSE) stream
- `GET /api/status` - Current solve status and result

**Concurrency Safety:**
- Mutex-protected `currentResult` for thread-safe access
- Context cancellation for request cancellation
- New solver instance per solve to prevent state leakage

**SSE Implementation:**
- Uses HTTP `text/event-stream` content type
- HTMX-compatible format: `data: {json}\n\n`
- 30-second timeout to prevent hanging connections
- Proper flushing for real-time updates

### 3. Frontend (HTMX)

**Architecture:**
- **EventSource API**: Native browser SSE support (no HTMX SSE extension needed)
- **Queue-based Rendering**: Moves queued during solving
- **Animated Playback**: Renders moves sequentially (0.5s intervals) after solution
- **State Reset**: Clears board and queues on reset/new solve

**Rendering Strategy:**
1. During solving: Collect moves in queue (no rendering)
2. On completion: Fetch final moves from `/api/status`
3. Animated playback: Render each move every 500ms
4. Visual feedback: Current move highlighted (pulse animation)

## Concurrency Patterns

### Goroutine Management

```go
// Solver runs in background goroutine
wg.Add(1)
go func() {
    defer wg.Done()
    success = s.solveRecursive(ctx, b, startPos, 1)
    // Signal completion via doneChan
}()
```

### Channel Communication

```go
// Non-blocking send with context awareness
select {
case s.moveChan <- update:
    // Successfully sent
case <-ctx.Done():
    return false // Cancelled
}
```

### WaitGroup Usage

- Ensures solver goroutine completes before returning result
- Prevents race conditions on result access
- Enables proper cleanup

## Memory Management

### Move Storage Strategy

**During Solving:**
- Moves appended to slice in real-time
- Backtracking removes moves from slice
- Buffer size grows/shrinks dynamically

**After Solving:**
- **Success**: Moves copied (deep copy) to result, original cleared
- **Failure**: Moves slice truncated to zero length, channels drained
- Prevents memory leaks from abandoned solves

### Channel Drainage

```go
func (s *Solver) clearChannels() {
    // Drain move channel
    for {
        select {
        case <-s.moveChan:
        default:
            goto doneMoves
        }
    }
    doneMoves:
    // Drain done channel
    select {
    case <-s.doneChan:
    default:
    }
}
```

## Production Considerations

### Scalability

**Current Design:**
- Single solver instance (sequential solves)
- In-memory state storage
- Suitable for low-to-medium traffic

**Future Enhancements:**
- Worker pool for parallel solves
- Redis for distributed state management
- WebSocket for bidirectional communication
- Rate limiting per IP/client

### Observability

**Missing (Production Additions Needed):**
- Structured logging (zap/logrus)
- Metrics (Prometheus)
- Distributed tracing (OpenTelemetry)
- Health check endpoint (`/health`)

### Error Handling

**Current:**
- Context cancellation handled gracefully
- HTTP errors return proper status codes
- Logging for debugging

**Production Improvements:**
- Retry logic for transient failures
- Circuit breaker pattern
- Error recovery strategies
- User-friendly error messages

### Security

**Current:**
- No authentication/authorization
- CORS headers for development
- Input validation on board size

**Production Requirements:**
- Request rate limiting
- Input sanitization
- CSRF protection
- Secure headers (CSP, HSTS)

## Testing Strategy

**Recommended Tests:**

1. **Unit Tests:**
   - Board operations (valid move, complete check)
   - Solver algorithm correctness
   - Channel communication

2. **Integration Tests:**
   - HTTP endpoint responses
   - SSE stream delivery
   - End-to-end solve flow

3. **Concurrency Tests:**
   - Multiple simultaneous solves
   - Channel buffer overflow
   - Context cancellation

4. **Performance Tests:**
   - Solve time benchmarks
   - Memory usage profiling
   - Concurrent request handling

## Performance Characteristics

**Expected Metrics (8x8 board):**
- **Solve Time**: < 1 second (Warnsdorff's heuristic)
- **Attempt Count**: 64-500 recursive calls
- **Memory**: ~10KB for move storage
- **Concurrent Requests**: Limited by single solver instance

**Bottlenecks:**
- Sequential solve operations
- Channel buffer size (may need tuning)
- Frontend rendering rate (500ms per move)

## Deployment

**Current:**
- Single binary (`go build ./cmd/server`)
- Static file serving from `web/`
- No external dependencies

**Production Deployment:**
- Docker containerization
- Reverse proxy (nginx/traefik)
- Process manager (systemd/supervisor)
- Health checks for orchestration
- Graceful shutdown (SIGTERM handling)

## Future Enhancements

1. **Multiple Algorithms**: Plugin architecture for different solving strategies
2. **Configuration**: YAML/JSON config for board size, algorithm, timeouts
3. **Persistence**: Save/load solutions to database
4. **Analytics**: Track solve times, success rates, popular starting positions
5. **Real-time Collaboration**: WebSocket for multi-user solving
6. **Optimization**: Parallel path exploration, memoization, A* variants

