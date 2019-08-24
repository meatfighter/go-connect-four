package main

// TODO REMOVE:
// go install github.com/meatfighter/go-connect-four/main

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	boardWidth  = 7
	boardHeight = 6

	cellEmpty = 0
	cellBlack = 1
	cellWhite = -1
)

var randomBits = func() [][][]uint64 {
	a := [][][]uint64{newUInt642D(boardWidth, boardHeight), newUInt642D(boardWidth, boardHeight)}
	for c := 1; c >= 0; c-- {
		for y := boardHeight - 1; y >= 0; y-- {
			for x := boardWidth - 1; x >= 0; x-- {
				a[c][y][x] = rand.Uint64()
			}
		}
	}
	return a
}()

func printRandomBits() {
	for c := 1; c >= 0; c-- {
		for y := boardHeight - 1; y >= 0; y-- {
			for x := boardWidth - 1; x >= 0; x-- {
				fmt.Printf("%d %d %d %d\n", c, y, x, randomBits[c][y][x])
			}
		}
	}
}

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
	board     [][]int
	ys        []int
	hash      uint64
	heuristic int
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
			switch row[x] {
			case cellEmpty:
				sb.WriteString(". ")
			case cellBlack:
				sb.WriteString("X ")
			case cellWhite:
				sb.WriteString("O ")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func newUInt642D(width, height int) [][]uint64 {
	a := make([][]uint64, height)
	rows := make([]uint64, width*height)
	for y := height - 1; y >= 0; y-- {
		a[y] = rows[width*y : width*(y+1)]
	}
	return a
}

func newInt2D(width, height int) [][]int {
	a := make([][]int, height)
	rows := make([]int, width*height)
	for y := height - 1; y >= 0; y-- {
		a[y] = rows[width*y : width*(y+1)]
	}
	return a
}

func newNode() *node {
	n := &node{}
	n.ys = make([]int, boardWidth)
	n.board = newInt2D(boardWidth, boardHeight)
	n.reset()
	return n
}

func (n *node) reset() {
	n.hash = 0
	n.heuristic = 0
	maxHeight := boardHeight - 1
	for x := boardWidth - 1; x >= 0; x-- {
		n.ys[x] = maxHeight
	}
	for y := maxHeight; y >= 0; y-- {
		row := n.board[y]
		for x := boardWidth - 1; x >= 0; x-- {
			row[x] = cellEmpty
		}
	}
}

func (n *node) isValidMove(x int) bool {
	return x >= 0 && x < boardWidth && n.ys[x] >= 0
}

func (n *node) makeMove(x, cell int) {
	y := n.ys[x]
	n.ys[x]--
	n.board[y][x] = cell
	if cell == cellBlack {
		n.heuristic += waysToWin[y][x]
		n.hash ^= randomBits[0][y][x]
	} else {
		n.heuristic -= waysToWin[y][x]
		n.hash ^= randomBits[1][y][x]
	}
}

func (n *node) undoMove(x, cell int) {
	n.ys[x]++
	y := n.ys[x]
	n.board[y][x] = cellEmpty
	if cell == cellBlack {
		n.heuristic -= waysToWin[y][x]
		n.hash ^= randomBits[0][y][x]
	} else {
		n.heuristic += waysToWin[y][x]
		n.hash ^= randomBits[1][y][x]
	}
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

}
