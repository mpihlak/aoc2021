package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MostCommonBit(codes []string, pos int) rune {
	countOnes := 0
	for _, code := range codes {
		c := code
		r := []rune(c)
		if r[pos] == '1' {
			countOnes += 1
		}
	}
	countZeroes := len(codes) - countOnes
	if countOnes >= countZeroes {
		return '1'
	} else {
		return '0'
	}
}

func CalculateRating(codes []string, mostCommon bool) string {
	positions := len(codes[0])

	for pos := 0; pos < positions; pos++ {
		bit := MostCommonBit(codes, pos)

		newCodes := []string{}
		for _, code := range codes {
			r := []rune(code)
			if mostCommon {
				if r[pos] == bit {
					newCodes = append(newCodes, code)
				}
			} else {
				if r[pos] != bit {
					newCodes = append(newCodes, code)
				}
			}
		}
		if len(newCodes) < 1 {
			break
		}
		codes = newCodes
	}

	return codes[0]
}

func main() {
	buf, _ := os.ReadFile("input.txt")
	input := string(buf)

	codes := make([]string, 0)
	for _, line := range strings.Split(input, "\n") {
		if line != "" {
			codes = append(codes, line)
		}
	}

	oneCounts := make([]int, len(codes[0]))
	for _, code := range codes {
		for pos, c := range code {
			if c == '1' {
				oneCounts[pos] += 1
			}
		}
	}
	gamma := ""
	epsilon := ""
	half := len(codes) / 2
	for _, count := range oneCounts {
		if count > half {
			gamma = gamma + "1"
			epsilon = epsilon + "0"
		} else {
			gamma = gamma + "0"
			epsilon = epsilon + "1"
		}
	}

	gammaVal, _ := strconv.ParseInt(gamma, 2, 64)
	epsilonVal, _ := strconv.ParseInt(epsilon, 2, 64)
	fmt.Println("Gamma =", gamma, " Epsilon=", epsilon)
	fmt.Println("GammaVal =", gammaVal, " EpsilonVal =", epsilonVal, " Answer =", gammaVal*epsilonVal)

	fmt.Println(oneCounts)
	O2Val, _ := strconv.ParseInt(CalculateRating(codes, true), 2, 64)
	CO2Val, _ := strconv.ParseInt(CalculateRating(codes, false), 2, 64)
	fmt.Println("O2 =", O2Val, "CO2 =", CO2Val, "Answer =", O2Val*CO2Val)
}
