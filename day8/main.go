package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Count the number of characters in s1 contained in s2
func ContainsChars(s1, s2 string) int {
	count := 0
	for _, c := range s1 {
		if strings.ContainsRune(s2, c) {
			count += 1
		}
	}
	return count
}

// Sort the chars in the string to make different permutations comparable
func Normalize(s string) string {
	chars := []byte(s)
	sort.Slice(chars, func(a, b int) bool { return chars[a] < chars[b] })
	return string(chars)
}

func Val(digits []string, mapSegmentsToValue map[string]int) int {
	val := 0
	for _, digit := range digits {
		val = val*10 + mapSegmentsToValue[Normalize(digit)]
	}
	return val
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	simpleCount := 0
	sum := 0

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		leftRight := strings.Split(line, " | ")
		left := leftRight[0]
		right := leftRight[1]
		rightDigits := []string{}

		// Part 1 - count the digits 1,4,7 and 8 on the right side
		for _, digit := range strings.Split(right, " ") {
			l := len(digit)
			if l == 2 || l == 4 || l == 3 || l == 7 {
				simpleCount += 1
			}
			rightDigits = append(rightDigits, digit)
		}

		digits := []string{}

		// Note the values for digits 1 and 4, they'll become handy for decoding some other digits
		seg1 := ""
		seg4 := ""
		for _, digit := range strings.Split(left, " ") {
			segments := Normalize(digit)
			switch len(digit) {
			case 2:
				seg1 = segments
			case 4:
				seg4 = segments
			}
			digits = append(digits, digit)
		}

		mapSegmentsToValue := make(map[string]int)

		for _, digit := range digits {
			segments := Normalize(digit)
			val := 0
			switch len(segments) {
			case 2:
				val = 1
			case 3:
				val = 7
			case 4:
				val = 4
			case 5:
				// either 5, 2 or 3
				if ContainsChars(seg1, segments) == 2 {
					// 3 if it contains 1
					val = 3
				} else {
					if ContainsChars(seg4, segments) == 3 {
						val = 5
					} else {
						val = 2
					}
				}
			case 6:
				// either 0, 9 or 6
				// 9 if it contains 4
				// else 0 if it contains 1
				// otherwise 6
				if ContainsChars(seg4, segments) == 4 {
					val = 9
				} else if ContainsChars(seg1, segments) == 2 {
					val = 0
				} else {
					val = 6
				}
			case 7:
				val = 8
			}
			mapSegmentsToValue[segments] = val
		}

		val := Val(rightDigits, mapSegmentsToValue)
		sum += val
	}

	fmt.Println("Count of 1,4,7,8 =", simpleCount)
	fmt.Println("Sum =", sum)
}
