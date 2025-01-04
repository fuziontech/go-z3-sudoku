package sudoku

import (
	"fmt"
	"strings"

	"github.com/fuziontech/go-z3"
)

// BoardView represents the visual state of the board
type BoardView map[int]map[int]CellView

// CellView represents a single cell's display properties
type CellView struct {
	value     byte
	fromInput bool
}

// ANSI color codes
const (
	Clear   = "\x1b[0;0m"
	RedFg   = "\x1b[31m"
	GreenFg = "\x1b[32m"
	Yellow  = "\x1b[1;43m"
	Green   = "\x1b[1;42m"
	BlueFg  = "\x1b[0;44m"
	White   = "\x1b[1;37m"
)

// Box drawing characters
const (
	TopRightCorner    = "\x1b(0\x6b"
	TopLeftCorner     = "\x1b(0\x6c"
	RightCorner       = "\x1b(0\x75"
	LeftCorner        = "\x1b(0\x74"
	BottomRightCorner = "\x1b(0\x6a"
	BottomLeftCorner  = "\x1b(0\x6d"
	Vertical          = "\x1b(0\x78\x1b(0"
	Horizontal        = "\x1b(0\x71"
	TopSplit          = "\x1b(0\x77"
	Split             = "\x1b(0\x6e"
	BottomSplit       = "\x1b(0\x76"
	EOL               = "\x1b(B\n"
)

// FromInput creates a string representation of the input board
func FromInput(input []byte) string {
	view := viewFromInput(input)
	return render(view)
}

// FromModel creates a string representation of the solved board
func FromModel(model *Model, solution *z3.Model) string {
	view := viewFromModel(model, solution)
	return render(view)
}

// render creates the string representation of the board with borders and colors
func render(board BoardView) string {
	var out strings.Builder

	out.WriteString(horizontalTop())
	for y := 0; y < BoardSize; y++ {
		if y == 3 || y == 6 {
			out.WriteString(horizontal())
		}
		for x := 0; x < BoardSize; x++ {
			if x%3 == 0 {
				out.WriteString(Vertical)
			}

			cell := board[x][y]
			if cell.value == 0 {
				out.WriteString(fmt.Sprintf("%s . %s", BlueFg, Clear))
			} else {
				if cell.fromInput {
					out.WriteString(fmt.Sprintf("%s %d %s", Yellow, cell.value, Clear))
				} else {
					out.WriteString(fmt.Sprintf("%s %d %s", Green, cell.value, Clear))
				}
			}
		}
		out.WriteString(Vertical)
		out.WriteString("\n")
	}
	out.WriteString(horizontalBottom())

	return out.String()
}

// horizontalTop creates the top border of the board
func horizontalTop() string {
	parts := []string{
		TopLeftCorner,
		strings.Repeat(Horizontal, 9),
		TopSplit,
		strings.Repeat(Horizontal, 9),
		TopSplit,
		strings.Repeat(Horizontal, 9),
		TopRightCorner,
		EOL,
	}
	return fmt.Sprintf("%s%s%s", White, strings.Join(parts, ""), Clear)
}

// horizontal creates the middle horizontal borders
func horizontal() string {
	parts := []string{
		LeftCorner,
		strings.Repeat(Horizontal, 9),
		Split,
		strings.Repeat(Horizontal, 9),
		Split,
		strings.Repeat(Horizontal, 9),
		RightCorner,
		EOL,
	}
	return fmt.Sprintf("%s%s%s", White, strings.Join(parts, ""), Clear)
}

// horizontalBottom creates the bottom border of the board
func horizontalBottom() string {
	parts := []string{
		BottomLeftCorner,
		strings.Repeat(Horizontal, 9),
		BottomSplit,
		strings.Repeat(Horizontal, 9),
		BottomSplit,
		strings.Repeat(Horizontal, 9),
		BottomRightCorner,
		EOL,
	}
	return fmt.Sprintf("%s%s%s", White, strings.Join(parts, ""), Clear)
}

// viewFromInput creates a BoardView from raw input bytes
func viewFromInput(input []byte) BoardView {
	out := make(BoardView)
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			if _, exists := out[x]; !exists {
				out[x] = make(map[int]CellView)
			}

			ch := input[y*BoardSize+x]
			var cell CellView
			if ch == '.' {
				cell = CellView{value: 0, fromInput: false}
			} else {
				value := ch - '0'
				if value < 1 || value > 9 {
					panic(fmt.Sprintf("illegal input value: %d must be in range [1..9]", value))
				}
				cell = CellView{value: value, fromInput: true}
			}
			out[x][y] = cell
		}
	}
	return out
}

// viewFromModel creates a BoardView from a solved model
func viewFromModel(model *Model, solution *z3.Model) BoardView {
	out := make(BoardView)
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			if _, exists := out[x]; !exists {
				out[x] = make(map[int]CellView)
			}

			pos := NewPos(x, y)
			cell := model.board[pos]
			solved := solution.Eval(cell.GetValue())
			value := byte(solved.Int())
			out[x][y] = CellView{
				value:     value,
				fromInput: cell.IsFromInput(),
			}
		}
	}
	return out
}
