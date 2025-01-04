package sudoku

import (
	"fmt"

	"github.com/fuziontech/go-z3"
)

const BoardSize = 9

// Pos represents a position on the Sudoku board
type Pos struct {
	x int
	y int
}

// NewPos creates a new position
func NewPos(x, y int) Pos {
	return Pos{x: x, y: y}
}

// Cell represents a cell in the Sudoku board
type Cell struct {
	value     *z3.AST
	fromInput bool
}

// NewDefaultCell creates a new cell with a fresh Z3 constant
func NewDefaultCell(ctx *z3.Context, pos Pos) *Cell {
	return &Cell{
		value:     ctx.FreshInt(fmt.Sprintf("value_x%d_y%d", pos.x, pos.y)),
		fromInput: false,
	}
}

// NewInputCell creates a new cell with a fixed input value
func NewInputCell(ctx *z3.Context, value int64) *Cell {
	return &Cell{
		value:     ctx.Int(int(value), ctx.IntSort()),
		fromInput: true,
	}
}

// GetValue returns the Z3 AST representing the cell's value
func (c *Cell) GetValue() *z3.AST {
	return c.value
}

// IsFromInput returns whether this cell was part of the input
func (c *Cell) IsFromInput() bool {
	return c.fromInput
}

// applyConstraints adds the basic constraints for this cell to the solver
func (c *Cell) applyConstraints(ctx *z3.Context, solver *z3.Solver) {
	one := ctx.Int(1, ctx.IntSort())
	nine := ctx.Int(9, ctx.IntSort())

	solver.Assert(c.value.Ge(one))
	solver.Assert(c.value.Lt(nine))
}

// Model represents the entire Sudoku board
type Model struct {
	ctx   *z3.Context
	board map[Pos]*Cell
}

// NewModel creates a new Sudoku model from input
func NewModel(ctx *z3.Context, input []byte) *Model {
	board := make(map[Pos]*Cell)

	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			pos := Pos{x: x, y: y}
			nextChar := input[y*BoardSize+x]

			if nextChar == '.' {
				board[pos] = NewDefaultCell(ctx, pos)
			} else {
				value := int64(nextChar - '0')
				board[pos] = NewInputCell(ctx, value)
			}
		}
	}

	return &Model{
		ctx:   ctx,
		board: board,
	}
}

// Solve attempts to solve the Sudoku puzzle
func Solve(input *Model) *z3.Model {
	solver := input.ctx.NewSolver()
	input.applyConstraints(solver)

	if solver.Check() != z3.True {
		return nil
	}

	return solver.Model()
}

// applyConstraints adds all Sudoku constraints to the solver
func (m *Model) applyConstraints(solver *z3.Solver) {
	// Apply basic cell constraints
	for _, cell := range m.board {
		cell.applyConstraints(m.ctx, solver)
	}

	// Apply row constraints
	for y := 0; y < BoardSize; y++ {
		rowCells := m.getRow(y)
		m.constrainDistinctValues(rowCells, solver)
	}

	// Apply column constraints
	for x := 0; x < BoardSize; x++ {
		colCells := m.getColumn(x)
		m.constrainDistinctValues(colCells, solver)
	}

	// Apply 3x3 cube constraints
	cubePositions := []Pos{
		{0, 0}, {0, 3}, {0, 6},
		{3, 0}, {3, 3}, {3, 6},
		{6, 0}, {6, 3}, {6, 6},
	}

	for _, pos := range cubePositions {
		cubeCells := m.getCube(pos)
		m.constrainDistinctValues(cubeCells, solver)
	}
}

// constrainDistinctValues adds a distinct constraint for a set of cells
func (m *Model) constrainDistinctValues(cells []*z3.AST, solver *z3.Solver) {
	solver.Assert(cells[0].Distinct(cells[1:]...))
}

// getCube returns all cell values in a 3x3 cube
func (m *Model) getCube(topLeft Pos) []*z3.AST {
	var cube []*z3.AST
	for y := topLeft.y; y < topLeft.y+3; y++ {
		for x := topLeft.x; x < topLeft.x+3; x++ {
			pos := Pos{x: x, y: y}
			cube = append(cube, m.board[pos].value)
		}
	}
	return cube
}

// getRow returns all cell values in a row
func (m *Model) getRow(targetY int) []*z3.AST {
	var row []*z3.AST
	for x := 0; x < BoardSize; x++ {
		pos := Pos{x: x, y: targetY}
		row = append(row, m.board[pos].value)
	}
	return row
}

// getColumn returns all cell values in a column
func (m *Model) getColumn(targetX int) []*z3.AST {
	var col []*z3.AST
	for y := 0; y < BoardSize; y++ {
		pos := Pos{x: targetX, y: y}
		col = append(col, m.board[pos].value)
	}
	return col
}
