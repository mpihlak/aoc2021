package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func SlidingWindowIncreaseCount(m []int) int {
	prevSum := 0
	incs := 0
	for i := 0; i < len(m)-2; i++ {
		sum := m[i] + m[i+1] + m[i+2]
		if prevSum > 0 && sum > prevSum {
			incs += 1
		}
		prevSum = sum
	}
	return incs
}

func main() {
	buf, _ := os.ReadFile("input.txt")
	input := string(buf)

	prev := 0
	incs := 0
	measurements := []int{}
	for _, line := range strings.Split(input, "\n") {
		if val, err := strconv.Atoi(line); err == nil {
			if prev > 0 && val > prev {
				incs += 1
			}
			prev = val
			measurements = append(measurements, val)
		}
	}

	fmt.Println("Increasing values =", incs)
	fmt.Println("Sliding window increasing values =", SlidingWindowIncreaseCount(measurements))
}
