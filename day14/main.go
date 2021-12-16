package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func InsertPolymers(template string, insertionRules map[string]string) string {
	result := ""
	runes := []rune(template)
	for i := 0; i < len(runes)-1; i++ {
		result += string(runes[i])
		result += insertionRules[template[i:i+2]]
	}
	result += template[len(template)-1 : len(template)]
	return result
}

func MinMaxQuantities(polymer string) (int, int) {
	freq := make(map[string]int)

	for _, c := range polymer {
		k := string(c)
		freq[k] += 1
	}

	min := math.MaxInt
	max := 0
	for _, v := range freq {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}

type CacheKey struct {
	a, b  string
	depth int
}

type FreqMap map[string]int64

type Result struct {
	a, b, c    string
	av, bv, cv int64
}

var frequencies = make(FreqMap)
var polymerCache = make(map[CacheKey]FreqMap)

func copyMap(src FreqMap) FreqMap {
	dst := make(FreqMap)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func FindPolymer(a, b string, insertionRules map[string]string, depth int, addB bool) FreqMap {
	if depth < 0 {
		return nil
	}

	cacheKey := CacheKey{a, b, depth}

	if cachedVal, ok := polymerCache[cacheKey]; ok {
		return cachedVal
	}

	c := insertionRules[a+b]

	if depth == 0 {
		bv := int64(0)
		if addB {
			bv = 1
		}
		res := make(FreqMap)
		res[a] += 1
		res[b] += bv
		return res
	}

	r1 := FindPolymer(a, c, insertionRules, depth-1, addB)
	r2 := FindPolymer(c, b, insertionRules, depth-1, addB)

	res := make(FreqMap)
	for k, v := range r1 {
		res[k] += v
	}
	for k, v := range r2 {
		res[k] += v
	}

	polymerCache[cacheKey] = res

	return res
}

func Polymerize(polymer string, insertionRules map[string]string, depth int) int64 {
	runes := []rune(polymer)
	for i := 0; i < len(runes)-1; i++ {
		a := string(runes[i])
		b := string(runes[i+1])
		m := FindPolymer(a, b, insertionRules, depth, false)
		for k, v := range m {
			frequencies[k] += v
		}
	}

	// Bump the last char
	frequencies[string(polymer[len(polymer)-1])] += 1

	values := []int64{}
	for _, v := range frequencies {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	return values[len(values)-1] - values[0]
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	template := ""
	insertionRules := make(map[string]string)

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		if template == "" {
			template = line
		} else {
			leftRight := strings.Split(line, " -> ")
			insertionRules[leftRight[0]] = leftRight[1]
		}
	}

	polymer := template
	for i := 0; i < 10; i++ {
		polymer = InsertPolymers(polymer, insertionRules)
		MinMaxQuantities(polymer)
	}

	min, max := MinMaxQuantities(polymer)
	fmt.Println("Part 1: min=", min, "max=", max, " diff=", max-min)

	result := Polymerize(template, insertionRules, 40)
	fmt.Println("Part 2: d=", result)
}
