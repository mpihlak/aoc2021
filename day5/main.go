package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Line struct {
	start Point
	end   Point
}

type OverlapMap map[Point]int

func (l Line) IsRightAngle() bool {
	return l.start.x == l.end.x || l.start.y == l.end.y
}

func calculateStep(x1, x2 int) int {
	if x1 == x2 {
		return 0
	} else if x1 < x2 {
		return 1
	} else {
		return -1
	}
}

func bump(x, y int, overlaps OverlapMap) {
	key := Point{x, y}
	newCount := 1
	if count, ok := overlaps[key]; ok {
		newCount += count
	}
	overlaps[key] = newCount
}

func WalkTheLine(l Line, overlaps OverlapMap) {
	xStep := calculateStep(l.start.x, l.end.x)
	yStep := calculateStep(l.start.y, l.end.y)
	x := l.start.x
	y := l.start.y
	for {
		bump(x, y, overlaps)

		if x == l.end.x && y == l.end.y {
			break
		}

		x += xStep
		y += yStep
	}
}

func CalculateOverlaps(lineList []Line, onlyRightAngles bool) int {
	overlapCount := make(OverlapMap)
	for _, line := range lineList {
		if onlyRightAngles && !line.IsRightAngle() {
			continue
		}
		WalkTheLine(line, overlapCount)
	}

	numOverlaps := 0
	for _, count := range overlapCount {
		if count > 1 {
			numOverlaps += 1
		}
	}

	return numOverlaps
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	lineList := []Line{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		endPointsStr := strings.Split(line, " -> ")
		xy := strings.Split(endPointsStr[0], ",")
		x1, _ := strconv.Atoi(xy[0])
		y1, _ := strconv.Atoi(xy[1])
		xy = strings.Split(endPointsStr[1], ",")
		x2, _ := strconv.Atoi(xy[0])
		y2, _ := strconv.Atoi(xy[1])
		lineList = append(lineList, Line{Point{x1, y1}, Point{x2, y2}})
	}

	overlappedPoints := CalculateOverlaps(lineList, true)
	fmt.Println("Overlapped points on horizontal and vertical lines =", overlappedPoints)
	overlappedPoints = CalculateOverlaps(lineList, false)
	fmt.Println("Overlapped points all lines =", overlappedPoints)
}
