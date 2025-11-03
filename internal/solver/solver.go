package solver

import (
	"context"
	"sync"
	"the_knight/pkg/board"
	"time"
)

// Solver manages the knight's tour solving algorithm with channel-based communication.
type Solver struct {
	mu sync.RWMutex
	// moveChan is buffered to prevent blocking the solver
	// Size 1000 handles rapid move sequences without significant delay
	moveChan chan MoveUpdate
	// doneChan signals completion (true = success, false = failure)
	doneChan chan bool
	// moves stores the sequence of moves (only if solution found)
	moves []MoveUpdate
	// attemptCount tracks recursive calls
	attemptCount int
}

// NewSolver creates a new solver instance with properly sized channels.
func NewSolver() *Solver {
	return &Solver{
		moveChan: make(chan MoveUpdate, 1000), // Buffered to prevent blocking
		doneChan: make(chan bool, 1),
		moves:    make([]MoveUpdate, 0, 64), // Pre-allocate for 8x8 board
	}
}

// Solve attempts to find a knight's tour solution using Warnsdorff's heuristic.
// It runs in a separate goroutine and communicates via channels.
func (s *Solver) Solve(ctx context.Context, boardSize int, startPos board.Position) (*SolveResult, error) {
	// Clear previous state
	s.mu.Lock()
	s.moves = s.moves[:0]
	s.attemptCount = 0
	s.mu.Unlock()

	// Drain channels to ensure clean state
	s.clearChannels()

	// Run solver in goroutine
	var wg sync.WaitGroup
	var success bool
	var solveErr error

	wg.Add(1)

	// Create a fresh board for this solve
	b := board.NewBoard(boardSize)

	go func() {
		defer wg.Done()
		success = s.solveRecursive(ctx, b, startPos, 1)
		// Signal completion (success or failure)
		// Note: solveRecursive sends doneChan internally when solution found,
		// but we need to ensure it's sent for failure case too
		if !success {
			select {
			case s.doneChan <- false:
			case <-ctx.Done():
			default:
				// Channel might be full or closed, ensure we signal somehow
				// Try one more time
				select {
				case s.doneChan <- false:
				case <-time.After(100 * time.Millisecond):
					// Give up if still can't send
				}
			}
		}
		// If success, doneChan was already sent by solveRecursive when solution found
	}()

	// Wait for completion or context cancellation
	select {
	case success = <-s.doneChan:
		// Solution found or failed - wait for goroutine to complete
		wg.Wait()
	case <-ctx.Done():
		// Context cancelled - clear channels and return
		solveErr = ctx.Err()
		s.clearChannels()
		wg.Wait()
		return &SolveResult{
			Success:      false,
			AttemptCount: s.getAttemptCount(),
		}, solveErr
	}

	// Only keep moves if solution was successful
	var finalMoves []MoveUpdate
	if success {
		s.mu.RLock()
		finalMoves = make([]MoveUpdate, len(s.moves))
		copy(finalMoves, s.moves)
		s.mu.RUnlock()
	} else {
		// Clear moves on failure - important for memory management
		s.mu.Lock()
		s.moves = s.moves[:0]
		s.mu.Unlock()
		s.clearChannels()
	}

	return &SolveResult{
		Success:      success,
		Moves:        finalMoves,
		AttemptCount: s.getAttemptCount(),
	}, nil
}

// solveRecursive implements the recursive backtracking algorithm with Warnsdorff's heuristic.
func (s *Solver) solveRecursive(ctx context.Context, b board.Board, currentPos board.Position, moveNumber int) bool {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return false
	default:
	}

	s.incAttemptCount()

	// Mark the current position
	b.WriteToBoard(currentPos, moveNumber)

	// Send move update (non-blocking with buffered channel)
	select {
	case s.moveChan <- MoveUpdate{Position: currentPos, MoveNumber: moveNumber, IsBacktrack: false}:
	case <-ctx.Done():
		return false
	}

	// Store move in sequence
	s.mu.Lock()
	s.moves = append(s.moves, MoveUpdate{Position: currentPos, MoveNumber: moveNumber, IsBacktrack: false})
	s.mu.Unlock()

	// Check if board is complete
	if b.IsComplete() {
		select {
		case s.doneChan <- true:
		case <-ctx.Done():
			return false
		}
		return true
	}

	// Knight move offsets (all 8 possible moves)
	knightMoves := []board.Position{
		{X: 2, Y: -1}, {X: 2, Y: 1}, {X: -2, Y: 1}, {X: -2, Y: -1},
		{X: 1, Y: 2}, {X: 1, Y: -2}, {X: -1, Y: 2}, {X: -1, Y: -2},
	}

	// Warnsdorff's heuristic: collect and sort by accessibility
	type MoveCandidate struct {
		position      board.Position
		accessibility int
	}

	var candidates []MoveCandidate

	for _, move := range knightMoves {
		newPos := board.Position{
			X: currentPos.X + move.X,
			Y: currentPos.Y + move.Y,
		}

		if b.IsValidMove(newPos) {
			accessibility := b.CountValidMoves(newPos)
			candidates = append(candidates, MoveCandidate{
				position:      newPos,
				accessibility: accessibility,
			})
		}
	}

	// Sort by accessibility (insertion sort for small lists)
	for i := 1; i < len(candidates); i++ {
		key := candidates[i]
		j := i - 1
		for j >= 0 && candidates[j].accessibility > key.accessibility {
			candidates[j+1] = candidates[j]
			j--
		}
		candidates[j+1] = key
	}

	// Try moves in priority order
	for _, candidate := range candidates {
		if s.solveRecursive(ctx, b, candidate.position, moveNumber+1) {
			return true
		}
	}

	// Backtrack: clear position and remove from moves
	b.ClearPosition(currentPos)

	// Send backtrack update
	select {
	case s.moveChan <- MoveUpdate{Position: currentPos, MoveNumber: 0, IsBacktrack: true}:
	case <-ctx.Done():
		return false
	}

	// Remove last move from sequence
	s.mu.Lock()
	if len(s.moves) > 0 {
		s.moves = s.moves[:len(s.moves)-1]
	}
	s.mu.Unlock()

	return false
}

// GetMoveChannel returns the channel for receiving move updates.
// Used by the web server to stream moves to clients.
func (s *Solver) GetMoveChannel() <-chan MoveUpdate {
	return s.moveChan
}

// clearChannels drains all channels to ensure clean state.
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

func (s *Solver) incAttemptCount() {
	s.mu.Lock()
	s.attemptCount++
	s.mu.Unlock()
}

func (s *Solver) getAttemptCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.attemptCount
}
