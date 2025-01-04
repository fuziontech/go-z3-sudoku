# Go Z3 Sudoku Solver

A Sudoku puzzle solver implemented in Go using Microsoft's Z3 theorem prover.

## Overview

This project demonstrates how to use Z3's constraint solving capabilities to solve Sudoku puzzles. It models the Sudoku rules as logical constraints and uses Z3 to find a solution that satisfies all the rules.

## Features

- Solve any valid 9x9 Sudoku puzzle
- Generate random puzzles with a specified number of starting values
- Load puzzles from input files
- Colorized output showing solution status

## Prerequisites

- Go 1.11 or higher
- Z3 theorem prover

## Building

1. Clone the repository
2. Build Z3 static library:
```
make
```

## Usage

```
go run .
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
