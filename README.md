# The Knight's Tour

A Go implementation of the Knight's Tour problem using Warnsdorff's heuristic algorithm with backtracking.

## Overview

The Knight's Tour is a classic chess problem: find a sequence of moves for a knight such that it visits every square on a chessboard exactly once. This implementation solves the problem using:

- **Warnsdorff's Heuristic**: Prioritizes moves that lead to positions with fewer available next moves
- **Backtracking**: Explores all possible paths and backtracks when no solution is found
- **Recursive Algorithm**: Efficiently searches the solution space

## Features

- ✅ Solves the Knight's Tour on an 8×8 chessboard
- ✅ Uses Warnsdorff's heuristic for fast solution finding
- ✅ Beautiful ASCII board visualization
- ✅ Execution timing and attempt counting
- ✅ Backtracking support for complete search

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
git clone <repository-url>
cd the_knight
```

### Running

```bash
go run main.go
```

Or build and run:

```bash
go build
./the_knight
```

## How It Works

1. **Initialization**: Creates an 8×8 board and starts at position (0, 0)
2. **Move Selection**: For each position, evaluates all 8 possible knight moves
3. **Heuristic Sorting**: Sorts valid moves by accessibility (fewest remaining moves first)
4. **Recursive Search**: Explores moves in priority order, marking visited squares
5. **Backtracking**: If a path fails, clears the position and tries the next option
6. **Completion Check**: Verifies when all squares have been visited

## Output

The program displays:
- Start timestamp
- End timestamp
- Total execution time
- Number of recursive attempts
- The solved board with move numbers

Example output:
```
Execution started at: 2024-01-15 10:30:45.123

Execution ended at: 2024-01-15 10:30:45.456

Time taken: 333ms
Total attempts: 64

Knight's tour completed!
+-----+-----+-----+-----+...
|  1  | 38  | 55  | ...
+-----+-----+-----+-----+...
...
```

## Algorithm Details

### Warnsdorff's Heuristic

The algorithm prioritizes moves that lead to positions with fewer available next moves. This strategy:
- Reduces the search space significantly
- Typically finds solutions in milliseconds instead of hours
- Minimizes backtracking

### Knight Moves

A knight can move to 8 positions from any square:
- 2 squares in one direction, 1 square perpendicular (±2, ±1) or (±1, ±2)

## Performance

- **Time Complexity**: O(8^N) worst case, but heuristic reduces this dramatically
- **Typical Execution**: < 1 second for 8×8 board
- **Attempt Count**: Usually 64-500 attempts (one per square with minimal backtracking)

## Project Structure

```
the_knight/
├── main.go      # Main program and board logic
├── go.mod       # Go module file
└── README.md    # This file
```

## License

This project is open source and available for educational purposes.
