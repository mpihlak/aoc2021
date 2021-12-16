package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CacheKey struct {
	timer int
	days  int
}

var fishCount = make(map[CacheKey]int64)

func CalculateLampfishCount(fishTimer, days int) int64 {
	if days < 0 {
		return 0
	}

	key := CacheKey{fishTimer, days}
	if val, ok := fishCount[key]; ok {
		return val
	}

	// Start with 1 to count self
	sum := int64(1)

	for days > 0 {
		days -= fishTimer + 1
		fishTimer = 6
		sum += CalculateLampfishCount(8, days)
	}

	fishCount[key] = sum

	return sum
}

func CalculateTotalCount(lampFishen []int, days int) int64 {
	count := int64(0)
	for _, fishTimer := range lampFishen {
		count += CalculateLampfishCount(fishTimer, days)
	}
	return count
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := strings.TrimSuffix(string(buf), "\n")

	lampFishen := []int{}
	for _, val := range strings.Split(input, ",") {
		timer, _ := strconv.Atoi(val)
		lampFishen = append(lampFishen, timer)
	}

	fmt.Println("Total fish after 80 days =", CalculateTotalCount(lampFishen, 80))
	fmt.Println("Total fish after 256 days =", CalculateTotalCount(lampFishen, 256))
}
