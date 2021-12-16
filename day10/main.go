package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func IsOpening(c rune) bool {
	return c == '(' || c == '[' || c == '{' || c == '<'
}

func MatchesOpening(opening, closing rune) bool {
	switch closing {
	case ')':
		return opening == '('
	case ']':
		return opening == '['
	case '}':
		return opening == '{'
	case '>':
		return opening == '<'
	}
	return false
}

func IsClosing(c rune) bool {
	return !IsOpening(c)
}

func UnmatchingScore(c rune) int {
	switch c {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	}
	return 0
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	unmatching := []rune{}
	completionScores := []int{}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		opening := []rune{}
		for _, c := range line {
			if IsOpening(c) {
				opening = append(opening, c)
			} else {
				oc := opening[len(opening)-1]
				opening = opening[:len(opening)-1]

				if !MatchesOpening(oc, c) {
					unmatching = append(unmatching, c)
					opening = []rune{} // Reset it, we won't need it
					break
				}
			}
		}

		score := 0
		if len(opening) > 0 {
			for i := range opening {
				oc := opening[len(opening)-1-i]
				values := map[rune]int{'(': 1, '[': 2, '{': 3, '<': 4}
				v := values[oc]
				score = score*5 + v
			}
			fmt.Printf("Unclosed line: %s: %s score: %d\n", line, string(opening), score)
			completionScores = append(completionScores, score)
		}
	}

	sum := 0
	for _, c := range unmatching {
		sum += UnmatchingScore(c)
	}
	fmt.Println("Syntax Error Score =", sum)

	sort.Ints(completionScores)
	middle := len(completionScores) / 2
	fmt.Println(completionScores, middle)
	fmt.Println("Completion Score =", completionScores[middle])
}
