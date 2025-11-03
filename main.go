package main

import (
	"fmt"
	"time"
)

var attemptCount int

func main() {
	var currentPosition Position = Position{0, 0}
	attemptCount = 0

	startTime := time.Now()
	fmt.Printf("Execution started at: %s\n\n", startTime.Format("2006-01-02 15:04:05.000"))

	board := NewBoard(8)
	success := board.MoveKnight(currentPosition, 1)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Execution ended at: %s\n\n", endTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Total attempts: %d\n\n", attemptCount)

	if success {
		fmt.Println("Knight's tour completed!")
		board.PrintBoard()
	} else {
		fmt.Println("No solution found for knight's tour")
		board.PrintBoard()
	}
}

type Board [][]int

type Position struct {
	x int
	y int
}

func NewBoard(size int) Board {
	board := make(Board, size, size)
	for i := range board {
		board[i] = make([]int, size, size)
	}
	return board
}

func (b Board) IsValidMove(pos Position) bool {
	return pos.x >= 0 && pos.x < len(b) && pos.y >= 0 && pos.y < len(b[0]) && b[pos.x][pos.y] == 0
}

// CountValidMoves counts how many valid moves are available from a given position
func (b Board) CountValidMoves(pos Position) int {
	count := 0
	validMoves := []Position{
		{2, -1}, {2, 1},
		{-2, 1}, {-2, -1},
		{1, 2}, {1, -2},
		{-1, 2}, {-1, -2},
	}

	for _, move := range validMoves {
		newPos := Position{
			x: pos.x + move.x,
			y: pos.y + move.y,
		}
		if b.IsValidMove(newPos) {
			count++
		}
	}
	return count
}

func (b Board) IsComplete() bool {
	for i := range b {
		for j := range b[i] {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (b Board) WriteToBoard(pos Position, moveNumber int) {
	b[pos.x][pos.y] = moveNumber
}

func (b Board) ClearPosition(pos Position) {
	b[pos.x][pos.y] = 0
}

func (b Board) MoveKnight(currentPosition Position, moveNumber int) bool {
	attemptCount++
	// Mark the current position
	b.WriteToBoard(currentPosition, moveNumber)

	// Check if board is complete
	if b.IsComplete() {
		return true
	}

	// Knight move offsets
	knightMoves := []Position{
		{2, -1}, {2, 1},
		{-2, 1}, {-2, -1},
		{1, 2}, {1, -2},
		{-1, 2}, {-1, -2},
	}

	// Warnsdorff's heuristic: collect valid moves with their accessibility counts
	type MoveCandidate struct {
		position      Position
		accessibility int // Number of valid moves from this position
	}

	var candidates []MoveCandidate

	// Collect all valid moves and their accessibility
	for _, move := range knightMoves {
		newPosition := Position{
			x: currentPosition.x + move.x,
			y: currentPosition.y + move.y,
		}

		if b.IsValidMove(newPosition) {
			// Count how many valid moves are available from this next position
			accessibility := b.CountValidMoves(newPosition)
			candidates = append(candidates, MoveCandidate{
				position:      newPosition,
				accessibility: accessibility,
			})
		}
	}

	// Sort candidates by accessibility (fewest moves available = highest priority)
	// We'll use a simple insertion sort for small lists (max 8 items)
	for i := 1; i < len(candidates); i++ {
		key := candidates[i]
		j := i - 1
		for j >= 0 && candidates[j].accessibility > key.accessibility {
			candidates[j+1] = candidates[j]
			j--
		}
		candidates[j+1] = key
	}

	// Try moves in order of least accessibility first (Warnsdorff's rule)
	for _, candidate := range candidates {
		if b.MoveKnight(candidate.position, moveNumber+1) {
			return true // Solution found
		}
	}

	// No valid path found, backtrack by clearing this position
	b.ClearPosition(currentPosition)
	return false
}

func (b Board) PrintBoard() {
	if len(b) == 0 {
		return
	}

	// Calculate the maximum width needed for numbers
	maxNum := len(b) * len(b[0])
	width := 0
	temp := maxNum
	for temp > 0 {
		temp /= 10
		width++
	}
	if width == 0 {
		width = 1
	}

	// Print top border
	fmt.Print("+")
	for j := 0; j < len(b[0]); j++ {
		for k := 0; k < width+2; k++ {
			fmt.Print("-")
		}
		fmt.Print("+")
	}
	fmt.Println()

	// Print each row
	for i := range b {
		fmt.Print("|")
		for j := range b[i] {
			fmt.Printf(" %*d |", width, b[i][j])
		}
		fmt.Println()

		// Print row separator
		fmt.Print("+")
		for j := 0; j < len(b[0]); j++ {
			for k := 0; k < width+2; k++ {
				fmt.Print("-")
			}
			fmt.Print("+")
		}
		fmt.Println()
	}
}
