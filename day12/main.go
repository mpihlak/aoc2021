package main

import (
	"fmt"
	"os"
	"strings"
)

type Graph map[string][]string
type VisitMap map[string]int

func addPair(caves Graph, a, b string) {
	if nodes, ok := caves[a]; ok {
		caves[a] = append(nodes, b)
	} else {
		caves[a] = []string{b}
	}
}

func isSmallCave(cave string) bool {
	return strings.ToLower(cave) == cave
}

func FindPaths(caves Graph, visited VisitMap, node string, count *int, maxVisits int, haveTwoVisits bool) {
	if node == "end" {
		*count += 1
		return
	}

	if isSmallCave(node) {
		if visitCount, ok := visited[node]; ok {
			if visitCount >= maxVisits {
				return
			}
			if visitCount >= 1 {
				if haveTwoVisits {
					return
				}
				haveTwoVisits = true
			}
		}
		visited[node] += 1
	}

	for _, caveToVisit := range caves[node] {
		FindPaths(caves, visited, caveToVisit, count, maxVisits, haveTwoVisits)
	}

	if isSmallCave(node) {
		visited[node] -= 1
	}
}

func FindAllPaths(caves Graph, maxSmallCaveVisits int) int {
	visited := make(VisitMap)
	visited["start"] = maxSmallCaveVisits
	count := 0

	for _, node := range caves["start"] {
		FindPaths(caves, visited, node, &count, maxSmallCaveVisits, false)
	}

	return count
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	caves := make(Graph)
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		leftRight := strings.Split(line, "-")
		start := leftRight[0]
		end := leftRight[1]

		addPair(caves, start, end)
		addPair(caves, end, start)
	}

	count := FindAllPaths(caves, 1)
	fmt.Println("Number of paths found (1 visit) =", count)
	count = FindAllPaths(caves, 2)
	fmt.Println("Number of paths found (2 visits) =", count)
}
