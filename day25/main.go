package main

import (
	"fmt"
	"os"
	"strings"
)

func printGrid(grid [][]rune) {
	for _, row := range grid {
		for _, c := range row {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	fmt.Println()
}

func moveEast(grid [][]rune) ([][]rune, bool) {
	result := make([][]rune, len(grid))
	moved := false
	for row := range grid {
		result[row] = make([]rune, len(grid[row]))
		for i := range result[row] {
			result[row][i] = '.'
		}
		for col, c := range grid[row] {
			newCol := (col + 1) % len(grid[row])
			if c == '>' && grid[row][newCol] == '.' {
				result[row][newCol] = c
				moved = true
			} else if result[row][col] == '.' {
				result[row][col] = c
			}
		}
	}
	return result, moved
}

func moveSouth(grid [][]rune) ([][]rune, bool) {
	result := make([][]rune, len(grid))
	moved := false
	for row := range grid {
		result[row] = make([]rune, len(grid[row]))
		for i := range result[row] {
			result[row][i] = '.'
		}
	}
	for row := range grid {
		for col, c := range grid[row] {
			newRow := (row + 1) % len(grid)
			if c == 'v' && grid[newRow][col] == '.' {
				result[newRow][col] = c
				moved = true
			} else if result[row][col] == '.' {
				result[row][col] = c
			}
		}
	}
	return result, moved
}

func moveCucumbers(grid [][]rune) ([][]rune, bool) {
	grid1, movedEast := moveEast(grid)
	grid2, movedSouth := moveSouth(grid1)
	return grid2, movedEast || movedSouth
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	grid := [][]rune{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		row := make([]rune, len(line))
		for pos, c := range line {
			row[pos] = c
		}

		grid = append(grid, row)
	}

	step := 1
	for {
		moved := false
		grid, moved = moveCucumbers(grid)
		if !moved {
			break
		}
		step++
	}

	fmt.Println("Stopped moving after", step, "steps.")
}
