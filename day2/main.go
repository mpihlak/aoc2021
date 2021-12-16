package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	buf, _ := os.ReadFile("input.txt")
	input := string(buf)

	depth := 0
	position := 0
	aim := 0
	aimedDepth := 0
	for _, line := range strings.Split(input, "\n") {
		instructions := strings.Split(line, " ")
		if len(instructions) == 2 {
			cmd, strVal := instructions[0], instructions[1]
			val, _ := strconv.Atoi(strVal)
			switch cmd {
			case "forward":
				position += val
				aimedDepth += val * aim
			case "down":
				depth += val
				aim += val
			case "up":
				depth -= val
				aim -= val
			}
		}
	}

	fmt.Printf("Part 1: Position = %v, Depth = %v, Answer = %v\n", position, depth, position*depth)
	fmt.Printf("Part 2: Position = %v, Depth = %v, Answer = %v\n", position, aimedDepth, position*aimedDepth)
}
