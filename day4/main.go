package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board [5][5]int
type BoardTags [5][5]bool

func NewBoard() Board {
	return [5][5]int{}
}

func (b Board) Print() {
	for _, row := range b {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (b Board) FindNumber(number int) (bool, int, int) {
	for row, cols := range b {
		for col, v := range cols {
			if v == number {
				return true, row, col
			}
		}
	}
	return false, -1, -1
}

func HaveBingoAfter(tags BoardTags, row, col int) bool {
	haveRow := true
	haveCol := true
	for i := 0; i < len(tags); i++ {
		if !tags[row][i] {
			haveRow = false
		}
		if !tags[i][col] {
			haveCol = false
		}
	}
	return haveRow || haveCol
}

func (b Board) SumUnmarkedNumbers(tags BoardTags) int {
	sum := 0
	for row, cols := range b {
		for col, v := range cols {
			if !tags[row][col] {
				sum += v
			}
		}
	}
	return sum
}

func PlayBingo(boardList []Board, numbersToDraw []int) (int, int) {
	boardTags := make([]BoardTags, len(boardList))
	hasWon := make([]bool, len(boardList))
	firstBoardScore := 0
	lastBoardScore := 0

	for _, number := range numbersToDraw {
		for boardIndex, board := range boardList {
			if found, row, col := board.FindNumber(number); found {
				boardTags[boardIndex][row][col] = true
				if HaveBingoAfter(boardTags[boardIndex], row, col) {
					score := board.SumUnmarkedNumbers(boardTags[boardIndex]) * number

					if firstBoardScore == 0 {
						firstBoardScore = score
					}
					if !hasWon[boardIndex] {
						lastBoardScore = score
						hasWon[boardIndex] = true
					}
				}
			}
		}
	}
	return firstBoardScore, lastBoardScore
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	numbersToDraw := []int{}
	boardList := []Board{}
	currentBoardRow := 0
	currentBoard := NewBoard()

	for n, line := range strings.Split(input, "\n") {
		if n == 0 {
			for _, val := range strings.Split(line, ",") {
				v, _ := strconv.Atoi(val)
				numbersToDraw = append(numbersToDraw, v)
			}
		} else if n == 1 {
			// Skip
		} else if line == "" {
			boardList = append(boardList, currentBoard)
			currentBoard = NewBoard()
			currentBoardRow = 0
		} else {
			for col := 0; col < 5; col++ {
				pos := col * 3
				s := strings.ReplaceAll(line[pos:pos+2], " ", "")
				if v, err := strconv.Atoi(s); err != nil {
					panic(line)
				} else {
					currentBoard[currentBoardRow][col] = v
				}
			}
			currentBoardRow += 1
		}
	}

	firstScore, lastScore := PlayBingo(boardList, numbersToDraw)
	fmt.Println("Score on first board to win =", firstScore)
	fmt.Println("Score on last board to win =", lastScore)
}
