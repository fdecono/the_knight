package solver

import "the_knight/pkg/board"

// MoveUpdate represents a single move in the knight's tour.
// Sent through channels to track progress in real-time.
type MoveUpdate struct {
	Position    board.Position
	MoveNumber  int
	IsBacktrack bool // true if this move is being backtracked (cleared)
}

// SolveResult encapsulates the result of a solve attempt.
type SolveResult struct {
	Success      bool
	Moves        []MoveUpdate
	AttemptCount int
}
