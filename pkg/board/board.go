package board

// Board represents a chess board as a 2D slice of integers.
// Each cell stores the move number (0 = unvisited).
type Board [][]int

// Position represents a coordinate on the board.
type Position struct {
	X int
	Y int
}

// NewBoard creates a new square board of the specified size, initialized with zeros.
func NewBoard(size int) Board {
	board := make(Board, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return board
}

// IsValidMove checks if a position is within bounds and unvisited.
func (b Board) IsValidMove(pos Position) bool {
	return pos.X >= 0 && pos.X < len(b) &&
		pos.Y >= 0 && pos.Y < len(b[0]) &&
		b[pos.X][pos.Y] == 0
}

// CountValidMoves returns the number of valid knight moves from a given position.
// This is used by Warnsdorff's heuristic.
func (b Board) CountValidMoves(pos Position) int {
	count := 0
	knightMoves := []Position{
		{2, -1}, {2, 1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}

	for _, move := range knightMoves {
		newPos := Position{X: pos.X + move.X, Y: pos.Y + move.Y}
		if b.IsValidMove(newPos) {
			count++
		}
	}
	return count
}

// IsComplete checks if all squares on the board have been visited.
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

// WriteToBoard marks a position with the given move number.
func (b Board) WriteToBoard(pos Position, moveNumber int) {
	b[pos.X][pos.Y] = moveNumber
}

// ClearPosition resets a position to unvisited (0).
func (b Board) ClearPosition(pos Position) {
	b[pos.X][pos.Y] = 0
}

// GetSize returns the board size (assuming square board).
func (b Board) GetSize() int {
	if len(b) == 0 {
		return 0
	}
	return len(b)
}

// GetCell returns the value at the specified position.
func (b Board) GetCell(pos Position) int {
	if pos.X < 0 || pos.X >= len(b) || pos.Y < 0 || pos.Y >= len(b[0]) {
		return -1
	}
	return b[pos.X][pos.Y]
}
