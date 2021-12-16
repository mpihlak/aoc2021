package main

import (
	"fmt"
	"os"
	"strings"
)

type Graph map[string][]string
type VisitMap map[string]int

type PathFinder struct {
	caves     Graph
	visited   VisitMap
	maxVisits int
}

func NewPathFinder(caves Graph, maxVisits int) PathFinder {
	p := PathFinder{caves, make(VisitMap), maxVisits}
	p.visited["start"] = maxVisits
	return p
}

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

func (p *PathFinder) FindPaths(node string, count *int, haveTwoVisits bool) {
	if node == "end" {
		*count += 1
		return
	}

	if isSmallCave(node) {
		if visitCount, ok := p.visited[node]; ok {
			if visitCount >= p.maxVisits {
				return
			}
			if visitCount >= 1 {
				if haveTwoVisits {
					return
				}
				haveTwoVisits = true
			}
		}
		p.visited[node] += 1
	}

	for _, caveToVisit := range p.caves[node] {
		p.FindPaths(caveToVisit, count, haveTwoVisits)
	}

	if isSmallCave(node) {
		p.visited[node] -= 1
	}
}

func (p *PathFinder) FindAllPaths() int {
	count := 0

	for _, node := range p.caves["start"] {
		p.FindPaths(node, &count, false)
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

	p1 := NewPathFinder(caves, 1)
	fmt.Println("Number of paths found (1 visit) =", p1.FindAllPaths())
	p2 := NewPathFinder(caves, 2)
	fmt.Println("Number of paths found (2 visits) =", p2.FindAllPaths())
}
