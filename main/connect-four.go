package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	boardWidth  = 7
	boardHeight = 6

	cellEmpty = 0
	cellBlack = 1
	cellWhite = -1

	ttExact      = 0
	ttLowerBound = 1
	ttUpperBound = 2

	nodeIntermediate = 0
	nodeBlackWins    = 1
	nodeWhiteWins    = 2
	nodeDraw         = 3

	infinity = 0x7FFFFFFF
)

var evaluationOrders = func() [][]int {
	middle := boardWidth >> 1
	orders := newInt2D(boardWidth, 1<<uint(middle))
	for i := len(orders) - 1; i >= 0; i-- {
		order := orders[i]
		var max int
		if (middle & 1) == 1 {
			order[boardWidth-1] = middle
			max = boardWidth - 1
		} else {
			max = boardWidth
		}
		b := i
		for j, k := 0, 0; j < max; j, k = j+2, k+1 {
			if (b & 1) == 0 {
				order[j] = k
				order[j+1] = boardWidth - 1 - k
			} else {
				order[j+1] = k
				order[j] = boardWidth - 1 - k
			}
			b >>= 1
		}
	}
	return orders
}()

func printEvaluationOrders() {
	for i := len(evaluationOrders) - 1; i >= 0; i-- {
		fmt.Println(evaluationOrders[i])
	}
}

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

func printWaysToWin() {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			fmt.Printf("%02d ", waysToWin[y][x])
		}
		fmt.Println()
	}
}

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
	nodeType  int
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
	n.nodeType = nodeIntermediate
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

// Capture nodeType before making move
func (n *node) makeMove(x, cell int) {

	y := n.ys[x]
	if y < 0 {
		return
	}

	n.ys[x]--
	n.board[y][x] = cell
	if cell == cellBlack {
		n.heuristic += waysToWin[y][x]
		n.hash ^= randomBits[0][y][x]
	} else {
		n.heuristic -= waysToWin[y][x]
		n.hash ^= randomBits[1][y][x]
	}

	var winner int
	if cell == cellBlack {
		winner = nodeBlackWins
	} else {
		winner = nodeWhiteWins
	}

	// Scan row for win
	row := n.board[y]
	p := x + 1
	for p < boardWidth && row[p] == cell {
		p++
	}
	q := x - 1
	for q >= 0 && row[q] == cell {
		q--
	}
	if p-q-1 >= 4 {
		n.nodeType = winner
		return
	}

	// Scan column for win
	p = y + 1
	for p < boardHeight && n.board[p][x] == cell {
		p++
	}
	q = y - 1
	for q >= 0 && n.board[q][x] == cell {
		q--
	}
	if p-q-1 >= 4 {
		n.nodeType = winner
		return
	}

	// Scan positive diagonal for win
	px := x + 1
	py := y + 1
	for px < boardWidth && py < boardHeight && n.board[py][px] == cell {
		px++
		py++
	}
	qx := x - 1
	qy := y - 1
	for qx >= 0 && qy >= 0 && n.board[qy][qx] == cell {
		qx--
		qy--
	}
	if px-qx-1 >= 4 {
		n.nodeType = winner
		return
	}

	// Scan negative diagonal for win
	px = x + 1
	py = y - 1
	for px < boardWidth && py >= 0 && n.board[py][px] == cell {
		px++
		py--
	}
	qx = x - 1
	qy = y + 1
	for qx >= 0 && qy < boardHeight && n.board[qy][qx] == cell {
		qx--
		qy++
	}
	if px-qx-1 >= 4 {
		n.nodeType = winner
		return
	}

	// Scan for draw
	for i := boardWidth - 1; i >= 0; i-- {
		if n.ys[i] >= 0 {
			n.nodeType = nodeIntermediate
			return
		}
	}
	n.nodeType = nodeDraw
}

// nodeType is type of node before move was made
func (n *node) undoMove(x, cell, nodeType int) {
	n.nodeType = nodeType
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

func createTTElement(flag, depth, value int) int32 {
	return int32((value << 8) | (depth << 2) | flag)
}

func extractTTValue(ttElement int32) int {
	return int(ttElement) >> 8
}

func extractTTDepth(ttElement int32) int {
	return (int(ttElement) >> 2) & 0x3F
}

func extractTTFlag(ttElement int32) int {
	return int(ttElement) & 0x03
}

func (n *node) computeMove(depth, alpha, beta, color int) int {

	move := -1
	nodeType := n.nodeType
	if depth == 0 || nodeType != nodeIntermediate {
		return move
	}

	tt := make(map[uint64]int32)
	order := evaluationOrders[rand.Int31n(int32(len(evaluationOrders)))]
	value := -infinity
	for i := len(order) - 1; i >= 0; i-- {
		m := order[i]
		if n.ys[m] >= 0 {
			n.makeMove(m, color)
			v := -n.negamax(depth-1, -beta, -alpha, -color, tt)
			n.undoMove(m, color, nodeType)
			if v > value {
				value = v
				move = m
			}
			if value > alpha {
				alpha = value
			}
			if alpha >= beta {
				break
			}
		}
	}

	return move
}

func (n *node) negamax(depth, alpha, beta, color int, tt map[uint64]int32) int {
	alphaOrig := alpha

	ttEntry, ttValidEntry := tt[n.hash]
	if ttValidEntry {
		ttEntryValue := extractTTValue(ttEntry)
		if extractTTDepth(ttEntry) >= depth {
			switch extractTTFlag(ttEntry) {
			case ttExact:
				return ttEntryValue
			case ttLowerBound:
				if ttEntryValue > alpha {
					alpha = ttEntryValue
				}
			case ttUpperBound:
				if ttEntryValue < beta {
					beta = ttEntryValue
				}
			}
		}
		if alpha >= beta {
			return ttEntryValue
		}
	}

	if depth == 0 {
		columnsRemaining := 0
		for x := boardWidth - 1; x >= 0 && columnsRemaining < 2; x-- {
			if n.ys[x] >= 0 {
				columnsRemaining++
			}
		}
		if columnsRemaining == 1 {
			depth = 1
		}
	}

	nodeType := n.nodeType
	if depth == 0 || nodeType != nodeIntermediate {
		value := n.heuristic
		if nodeType == nodeBlackWins {
			value += (depth + 1) << 12
		} else if nodeType == nodeWhiteWins {
			value -= (depth + 1) << 12
		}
		return color * value
	}

	order := evaluationOrders[rand.Int31n(int32(len(evaluationOrders)))]
	value := -infinity
	for i := len(order) - 1; i >= 0; i-- {
		move := order[i]
		if n.ys[move] >= 0 {
			n.makeMove(move, color)
			v := -n.negamax(depth-1, -beta, -alpha, -color, tt)
			n.undoMove(move, color, nodeType)
			if v > value {
				value = v
			}
			if value > alpha {
				alpha = value
			}
			if alpha >= beta {
				break
			}
		}
	}

	var ttFlag int
	if value <= alphaOrig {
		ttFlag = ttUpperBound
	} else if value >= beta {
		ttFlag = ttLowerBound
	} else {
		ttFlag = ttExact
	}
	tt[n.hash] = createTTElement(ttFlag, depth, value)

	return value
}

func prompt(reader *bufio.Reader, message string) string {
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	fmt.Println()
	return strings.TrimSpace(text)
}

func promptForValue(reader *bufio.Reader, message string) int {
	for {
		value, err := strconv.Atoi(prompt(reader, message))
		if err == nil {
			return value
		}
	}
}

func promptForValueWithDefault(reader *bufio.Reader, message string, defaultValue int) int {
	for {
		value, err := strconv.Atoi(prompt(reader, message))
		if err == nil {
			return value
		}
		return defaultValue
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	fmt.Println()
	fmt.Println("Connect Four: Human vs. Computer")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	maxDepth := 0
	for maxDepth < 4 || maxDepth > 20 {
		maxDepth = promptForValueWithDefault(reader, "Enter max search depth (4--20) [10]: ", 10)
	}

	for {
		n := newNode()
		var color int
		if (rand.Int() & 1) == 0 {
			color = cellBlack
			fmt.Println("Human (X) plays first.")
		} else {
			color = cellWhite
			fmt.Println("Computer (O) plays first.")
		}

		fmt.Println()
		n.print()

		for n.nodeType == nodeIntermediate {
			var move int
			if color == cellBlack {
				move = promptForValue(reader, "Enter column: ") - 1
			} else {
				move = n.computeMove(maxDepth, -infinity, infinity, color)
				fmt.Printf("Computer drops O into column %d.\n", move+1)
				fmt.Println()
			}
			if n.isValidMove(move) {
				n.makeMove(move, color)
				color = -color
			} else {
				fmt.Println("Invalid column.")
				fmt.Println()
			}

			n.print()
		}

		switch n.nodeType {
		case nodeBlackWins:
			fmt.Println("Human wins.")
		case nodeWhiteWins:
			fmt.Println("Computer wins.")
		case nodeDraw:
			fmt.Println("Draw.")
		}
		fmt.Println()

		s := prompt(reader, "Play again (y/n)? ")
		if len(s) == 0 || unicode.ToLower(rune(s[0])) != 'y' {
			break
		}
	}

	fmt.Println("Thanks for playing.")
	fmt.Println()
	fmt.Println("Goodbye.")
	fmt.Println()
}
