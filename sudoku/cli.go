package sudoku

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Args holds the command line arguments
type Args struct {
	Generate int
	Input    string
}

// NewArgs parses command line arguments and returns an Args struct
func NewArgs() *Args {
	args := &Args{}

	flag.IntVar(&args.Generate, "generate", 18, "Generate an input Sudoku board to solve; arg is number of seed values to apply")
	flag.IntVar(&args.Generate, "g", 18, "Generate an input Sudoku board to solve (shorthand)")

	flag.StringVar(&args.Input, "input", "", "Load Sudoku board from specified file (see 'example.board')")
	flag.StringVar(&args.Input, "i", "", "Load Sudoku board from specified file (shorthand)")

	flag.Parse()
	return args
}

// BoardFromFile reads a Sudoku board from the specified file
func (a *Args) BoardFromFile() ([]byte, error) {
	content, err := os.ReadFile(a.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to load file specified in --input arg: %v", err)
	}

	// Remove whitespace and convert to string
	board := strings.Map(func(r rune) rune {
		if !strings.ContainsRune(" \n\t\r", r) {
			return r
		}
		return -1
	}, string(content))

	if len(board) != 81 {
		return nil, fmt.Errorf("board must be exactly 81 characters (got %d)", len(board))
	}

	return []byte(board), nil
}

// GenerateBoard creates a new random Sudoku board with the specified number of seed values
func (a *Args) GenerateBoard() []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create empty board
	board := make([]byte, 81)
	for i := range board {
		board[i] = '.'
	}

	// Add random values
	for i := 0; i <= a.Generate; i++ {
		for {
			cellValue := byte(r.Intn(BoardSize) + 1 + '0')
			candidateX := r.Intn(BoardSize)
			candidateY := r.Intn(BoardSize)

			if a.validPlacement(board, candidateX, candidateY, cellValue) {
				board[(candidateY*BoardSize)+candidateX] = cellValue
				break
			}
		}
	}

	return board
}

// validPlacement checks if a value can be placed at the given position
func (a *Args) validPlacement(board []byte, x, y int, val byte) bool {
	// Check row
	for checkX := 0; checkX < BoardSize; checkX++ {
		if x != checkX {
			if board[(y*BoardSize)+checkX] == val {
				return false
			}
		}
	}

	// Check column
	for checkY := 0; checkY < BoardSize; checkY++ {
		if y != checkY {
			if board[(checkY*BoardSize)+x] == val {
				return false
			}
		}
	}

	// Check 3x3 cube
	top := (y / 3) * 3
	left := (x / 3) * 3
	for checkY := top; checkY < top+3; checkY++ {
		for checkX := left; checkX < left+3; checkX++ {
			if checkY != y && checkX != x {
				if board[(checkY*BoardSize)+checkX] == val {
					return false
				}
			}
		}
	}

	return true
}

// Example usage:
// func main() {
//     args := NewArgs()
//     var board []byte
//     var err error
//
//     if args.input != "" {
//         board, err = args.BoardFromFile()
//         if err != nil {
//             fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//             os.Exit(1)
//         }
//     } else {
//         board = args.GenerateBoard()
//     }
//
//     // Use the board...
// }
