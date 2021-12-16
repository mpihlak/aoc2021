package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fold(grid [][]int, x, y int) [][]int {
	width := len(grid[0])
	height := len(grid)
	if x > 0 {
		width = x
	}
	if y > 0 {
		height = y
	}

	newGrid := make([][]int, height)
	for i := 0; i < height; i++ {
		newGrid[i] = make([]int, width)
		copy(newGrid[i], grid[i])
	}

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if x > 0 {
				p := width + col + 1
				if p < len(grid[row]) {
					newGrid[row][width-col-1] += grid[row][width+col+1]
				}
			} else {
				p := height + row + 1
				if p < len(grid) {
					newGrid[height-row-1][col] += grid[p][col]
				}
			}
		}
	}

	return newGrid
}

func CountPoints(grid [][]int) int {
	count := 0
	for _, rows := range grid {
		for _, val := range rows {
			if val > 0 {
				count += 1
			}
		}
	}
	return count
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	pointCoords := make([][]int, 0)
	folds := make([][]int, 0)
	maxX := 0
	maxY := 0
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "fold") {
			s := strings.TrimPrefix(line, "fold along ")
			axisVal := strings.Split(s, "=")
			axis := axisVal[0]
			val, _ := strconv.Atoi(axisVal[1])
			if axis == "x" {
				folds = append(folds, []int{val, 0})
			} else {
				folds = append(folds, []int{0, val})
			}
		} else {
			xy := strings.Split(line, ",")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			pointCoords = append(pointCoords, []int{x, y})
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	grid := make([][]int, maxY+1)
	for y := 0; y <= maxY; y++ {
		grid[y] = make([]int, maxX+1)
	}

	for _, coords := range pointCoords {
		grid[coords[1]][coords[0]] = 1
	}

	for n, f := range folds {
		grid = Fold(grid, f[0], f[1])
		count := CountPoints(grid)
		if n == 0 {
			fmt.Println("Fold", n+1, "Points =", count)
		}
	}

	for _, row := range grid {
		for _, v := range row {
			if v > 0 {
				fmt.Printf("▇")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}

}
