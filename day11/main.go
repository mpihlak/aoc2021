package main

import (
	"fmt"
	"os"
	"strings"
)

type Point struct {
	row, col, val int
}

func ValidPos(g [][]int, row, col int) bool {
	return len(g) > 0 && row >= 0 && col >= 0 && row < len(g) && col < len(g[0])
}

func AddEnergy(g [][]int, row, col int) {
	val := &g[row][col]
	if *val > 9 {
		// Already flashed
		return
	}

	*val += 1

	if *val > 9 {
		// Flashed now, add energy to neighbors
		for r := row - 1; r <= row+1; r++ {
			for c := col - 1; c <= col+1; c++ {
				if ValidPos(g, r, c) && !(r == row && c == col) {
					AddEnergy(g, r, c)
				}
			}
		}
	}
}

func CalculateFlashes(grid [][]int) (int, int) {
	totalFlashes100 := 0
	whenAllFlashed := 0

	for step := 0; ; step++ {
		for rowNum, cols := range grid {
			for colNum := range cols {
				AddEnergy(grid, rowNum, colNum)
			}
		}

		allFlashed := true
		for r, cols := range grid {
			for c, val := range cols {
				if val > 9 {
					if step < 100 {
						totalFlashes100 += 1
					}
					grid[r][c] = 0
				} else {
					allFlashed = false
				}
			}
		}
		if allFlashed {
			whenAllFlashed = step + 1
			break
		}
	}

	return totalFlashes100, whenAllFlashed
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	grid := make([][]int, 0)
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		row := make([]int, len(line))
		for pos, c := range line {
			row[pos] = int(c - '0')
		}

		grid = append(grid, row)
	}

	flashes, whenAllFlashed := CalculateFlashes(grid)
	fmt.Println("Number of flashes =", flashes)
	fmt.Println("All octopuses flashed at", whenAllFlashed)
}
