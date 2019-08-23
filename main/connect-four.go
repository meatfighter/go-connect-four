package main

// TODO REMOVE:
// go install github.com/meatfighter/go-connect-four/main

import (
	"fmt"
	"strings"
)

const (
	boardWidth  = 7
	boardHeight = 6

	cellEmpty = 0
	cellX     = 1
	cellO     = 2
)

var cellNames = []rune(".XO")

var waysToWin = func() [][]int {
	w := newInt2D(boardWidth, boardHeight)
	for y := boardHeight - 1; y >= 0; y-- {
		for x := boardWidth - 1; x >= 0; x-- {
			w[y][x] = computeWaysToWin(w, x, y)
		}
	}
	return w
}()

func computeWaysToWin(w [][]int, x, y int) int {
	return computeWaysToWinWithDeltas(w, x, y, 1, 0) +
		computeWaysToWinWithDeltas(w, x, y, 0, 1) +
		computeWaysToWinWithDeltas(w, x, y, 1, 1) +
		computeWaysToWinWithDeltas(w, x, y, -1, 1)
}

func computeWaysToWinWithDeltas(w [][]int, x, y, dx, dy int) int {
	ways := computePositiveAdvance(w, x, y, dx, dy) + computePositiveAdvance(w, x, y, -dx, -dy) - 2
	if ways < 0 {
		return 0
	}
	return ways
}

func computePositiveAdvance(w [][]int, x, y, dx, dy int) int {
	p := 0
	for i, px, py := 1, x+dx, y+dy; i <= 3 && isValidCoordinate(px, py); i, px, py, p = i+1, px+dx, py+dy, p+1 {
	}
	return p
}

func isValidCoordinate(x, y int) bool {
	return x >= 0 && y >= 0 && x < boardWidth && y < boardHeight
}

type node struct {
	board   [][]int
	heights []int
}

func (n *node) print() {
	fmt.Println(convertBoardToString(n.board))
}

func convertBoardToString(board [][]int) string {
	var sb strings.Builder
	sb.WriteString("  ")
	for x := 0; x < boardWidth; x++ {
		sb.WriteRune(rune('1' + x))
		sb.WriteRune(' ')
	}
	sb.WriteRune('\n')
	for y := 0; y < boardHeight; y++ {
		row := board[y]
		sb.WriteRune(rune('1' + y))
		sb.WriteRune(' ')
		for x := 0; x < boardWidth; x++ {
			sb.WriteRune(cellNames[row[x]])
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func newInt2D(width, height int) [][]int {
	var a = make([][]int, height)
	rows := make([]int, width*height)
	for y := height - 1; y >= 0; y-- {
		a[y] = rows[width*y : width*(y+1)]
	}
	return a
}

func newNode() *node {
	n := &node{}
	n.heights = make([]int, boardWidth)
	n.board = newInt2D(boardWidth, boardHeight)
	return n
}

func printWaysToWin() {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			fmt.Printf("%02d ", waysToWin[y][x])
		}
		fmt.Println()
	}
}

func main() {
	printWaysToWin()
}
