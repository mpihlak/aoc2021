package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Point struct {
	row, col, val int
}

func ValidPos(g [][]int, row, col int) bool {
	return len(g) > 0 && row >= 0 && col >= 0 && row < len(g) && col < len(g[0])
}

func IsLowest(g [][]int, row, col int) bool {
	val := g[row][col]
	found := true
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if ValidPos(g, r, c) && !(r == row && c == col) {
				if g[r][c] <= val {
					found = false
				}
			}
		}
	}
	return found
}

func FindLowPoints(grid [][]int) []Point {
	lowPoints := []Point{}

	for rowNum, cols := range grid {
		for colNum, val := range cols {
			if IsLowest(grid, rowNum, colNum) {
				lowPoints = append(lowPoints, Point{rowNum, colNum, val})
			}
		}
	}
	return lowPoints
}

func CalculateRisk(points []Point) int {
	sum := 0
	for _, p := range points {
		sum += p.val + 1
	}
	return sum
}

func getBasinSize(g [][]int, val, row, col int, visited [][]bool) int {
	if g[row][col] < val || g[row][col] == 9 {
		return 0
	}

	visited[row][col] = true
	sum := 1
	neighbors := [][]int{{row - 1, col}, {row + 1, col}, {row, col - 1}, {row, col + 1}}
	for i := range neighbors {
		r := neighbors[i][0]
		c := neighbors[i][1]
		if ValidPos(g, r, c) && !visited[r][c] {
			sum += getBasinSize(g, g[row][col], r, c, visited)
		}
	}

	return sum
}

func CalculateBasinSize(g [][]int, p Point) int {
	visited := make([][]bool, len(g))
	for i := range g {
		visited[i] = make([]bool, len(g[i]))
	}

	return getBasinSize(g, p.val, p.row, p.col, visited)
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

	lowPoints := FindLowPoints(grid)
	fmt.Println("Sum of risk index =", CalculateRisk(lowPoints))

	basinSizes := []int{}
	for _, point := range lowPoints {
		size := CalculateBasinSize(grid, point)
		basinSizes = append(basinSizes, size)
	}

	sort.Ints(basinSizes)
	l := len(basinSizes)
	res := basinSizes[l-1] * basinSizes[l-2] * basinSizes[l-3]
	fmt.Println("Top 3 basin sizes multiplied =", res)
}
