package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func ValidPos(g [][]int, row, col int) bool {
	return len(g) > 0 && row >= 0 && col >= 0 && row < len(g) && col < len(g[0])
}

type Node struct {
	row, col, distance int
}

func FindLowestRiskPath(grid [][]int) int {
	distances := make([][]int, len(grid))
	for row := range grid {
		distances[row] = make([]int, len(grid[row]))
	}

	nodes := []Node{{0, 0, 0}}

	tryMove := func(row, col, d int) {
		if !ValidPos(grid, row, col) || distances[row][col] > 0 {
			return
		}

		destNode := Node{
			row, col, d + grid[row][col],
		}

		distances[row][col] = d + grid[row][col]
		nodes = append(nodes, destNode)
	}

	for len(nodes) > 0 {
		// Poor man's priority queue
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].distance < nodes[j].distance })

		n := nodes[0]
		nodes = nodes[1:]

		tryMove(n.row, n.col+1, n.distance)
		tryMove(n.row+1, n.col, n.distance)
		tryMove(n.row, n.col-1, n.distance)
		tryMove(n.row-1, n.col, n.distance)
	}

	endRow := len(grid) - 1
	endCol := len(grid[0]) - 1

	return distances[endRow][endCol]
}

func MultiplyGrid(grid [][]int, n int) [][]int {
	result := make([][]int, len(grid)*n)
	for row := range result {
		result[row] = make([]int, len(grid[0])*n)
	}

	for addRow := 0; addRow < n; addRow++ {
		for addCol := 0; addCol < n; addCol++ {
			r := addRow * len(grid)
			c := addCol * len(grid[0])

			for row := range grid {
				for col := range grid[row] {
					v := grid[row][col] + addRow + addCol
					if v > 9 {
						v = v%10 + 1
					}
					result[r+row][c+col] = v
				}
			}
		}
	}

	return result
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

	riskLevel := FindLowestRiskPath(grid)
	fmt.Println("Single grid risk level =", riskLevel)

	bigGrid := MultiplyGrid(grid, 5)

	riskLevel = FindLowestRiskPath(bigGrid)
	fmt.Println("Big grid risk level =", riskLevel)
}
