# Go Z3 Sudoku Solver

A Sudoku puzzle solver implemented in Go using Microsoft's Z3 theorem prover.

<img width="597" alt="image" src="https://github.com/user-attachments/assets/49bd2eba-8a71-4cdf-a506-9360cc20d0c6" />


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
2. Install Z3 on brew
```
brew install z3
```
3. Setup Z3 in your PATH
```
export CGO_CFLAGS="-I$(brew --prefix z3)/include"
export CGO_LDFLAGS="-L$(brew --prefix z3)/lib"
export DYLD_LIBRARY_PATH="$(brew --prefix z3)/lib"
```
4. Build Z3 static library:
```
make
```

## Usage

```
go run .
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
