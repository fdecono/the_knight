package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"the_knight/internal/solver"
	"the_knight/pkg/board"
)

// Server handles HTTP requests and manages the solver state.
type Server struct {
	solver        *solver.Solver
	mu            sync.RWMutex
	currentResult *solver.SolveResult
	templates     *template.Template
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewServer creates a new web server instance.
func NewServer() *Server {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		solver:    solver.NewSolver(),
		templates: tmpl,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start begins the HTTP server on the specified address.
func (s *Server) Start(addr string) error {
	// Serve static files if needed
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/solve", s.handleSolve)
	http.HandleFunc("/api/moves/stream", s.handleMoveStream)
	http.HandleFunc("/api/status", s.handleStatus)

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleIndex serves the main HTML page with HTMX.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleSolve starts a new solve operation.
func (s *Server) handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Cancel any existing solve and create new solver instance
	s.mu.Lock()
	if s.cancel != nil {
		s.cancel()
	}
	// Create new solver instance to reset state
	s.solver = solver.NewSolver()
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel
	s.currentResult = nil
	s.mu.Unlock()

	// Parse request
	var req struct {
		Size     int            `json:"size"`
		StartPos board.Position `json:"startPos"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Size <= 0 || req.Size > 20 {
		req.Size = 8 // Default to 8x8
	}

	// Start solving in background
	go func() {
		result, err := s.solver.Solve(ctx, req.Size, req.StartPos)
		if err != nil && err != context.Canceled {
			log.Printf("Solve error: %v", err)
			return
		}

		s.mu.Lock()
		if result != nil {
			s.currentResult = result
		}
		s.mu.Unlock()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "solving"})
}

// handleMoveStream streams moves via Server-Sent Events (SSE) for HTMX.
func (s *Server) handleMoveStream(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get move channel
	moveChan := s.solver.GetMoveChannel()

	// Flush headers
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// Stream moves
	for {
		select {
		case move := <-moveChan:
			// Send as HTMX SSE format
			data, _ := json.Marshal(move)
			fmt.Fprintf(w, "data: %s\n\n", string(data))

			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			// Check if we should stop (solution found or failed)
			s.mu.RLock()
			result := s.currentResult
			s.mu.RUnlock()

			if result != nil {
				if result.Success {
					// Send completion event
					fmt.Fprintf(w, "data: {\"type\":\"complete\",\"success\":true}\n\n")
				} else {
					fmt.Fprintf(w, "data: {\"type\":\"complete\",\"success\":false}\n\n")
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				return
			}

		case <-r.Context().Done():
			return
		case <-time.After(30 * time.Second):
			// Timeout to prevent hanging connections
			return
		}
	}
}

// handleStatus returns the current solve status.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	result := s.currentResult
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if result != nil {
		json.NewEncoder(w).Encode(result)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "not_started"})
	}
}
