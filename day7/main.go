package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func CalculateAlignmentCost(positions []int, adjustedFormula bool) int {
	minPos := math.MaxInt
	maxPos := 0
	for _, pos := range positions {
		if pos < minPos {
			minPos = pos
		}
		if pos > maxPos {
			maxPos = pos
		}
	}

	minSum := math.MaxInt
	for alignTo := minPos; alignTo <= maxPos; alignTo++ {
		sum := 0
		for _, pos := range positions {
			fuel := int(math.Abs(float64(pos - alignTo)))
			if adjustedFormula {
				fuel = fuel * (fuel + 1) / 2
			}
			sum += fuel
		}
		if sum < minSum {
			minSum = sum
		}
	}
	return minSum
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := strings.TrimSuffix(string(buf), "\n")

	positions := []int{}
	for _, strVal := range strings.Split(input, ",") {
		val, _ := strconv.Atoi(strVal)
		positions = append(positions, val)
	}

	fmt.Println("Part 1: Cost of alignment =", CalculateAlignmentCost(positions, false))
	fmt.Println("Part 2: Cost of alignment =", CalculateAlignmentCost(positions, true))
}
