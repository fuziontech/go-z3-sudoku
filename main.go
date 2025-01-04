package main

import (
	"fmt"
	"os"

	"github.com/fuziontech/go-z3"
	"github.com/fuziontech/go-z3-sudoku/sudoku"
)

func main() {
	args := sudoku.NewArgs()
	var input []byte
	var err error

	if args.Input != "" {
		fmt.Printf("Loading board from input file: %s\n", args.Input)
		input, err = args.BoardFromFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Generating random board with %d seed values\n", args.Generate)
		input = args.GenerateBoard()
	}

	fmt.Println("Input board:")
	fmt.Println(sudoku.FromInput(input))
	fmt.Println()

	// Initialize Z3
	cfg := z3.NewConfig()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()
	defer cfg.Close()

	// Create and solve the model
	board := sudoku.NewModel(ctx, input)
	solution := sudoku.Solve(board)

	if solution != nil {
		fmt.Printf("%sZ3 Solver result: SAT%s\n", sudoku.GreenFg, sudoku.Clear)
		fmt.Println(sudoku.FromModel(board, solution))
	} else {
		fmt.Printf("%sZ3 Solver result: UNSAT%s\n", sudoku.RedFg, sudoku.Clear)
	}
	fmt.Println()
}
