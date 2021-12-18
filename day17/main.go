package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MaxithinTarget(xs, ys, x1, x2, y1, y2 int) (int, bool) {
	x := 0
	y := 0
	maxY := 0
	found := false
	for {
		if y > maxY {
			maxY = y
		}

		if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
			found = true
		}

		x += xs
		y += ys
		if xs > 0 {
			xs -= 1
		} else if xs < 0 {
			xs += 1
		}

		ys -= 1

		if y2 < y1 && y < y2 {
			break
		}
		if y1 < y2 && y < y1 {
			break
		}
	}

	return maxY, found
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	var x1, x2, y1, y2 int

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "target area: ")
		lr := strings.Split(line, ", ")
		xr := strings.TrimPrefix(lr[0], "x=")
		yr := strings.TrimPrefix(lr[1], "y=")
		xvs := strings.Split(xr, "..")
		yvs := strings.Split(yr, "..")
		x1, _ = strconv.Atoi(xvs[0])
		x2, _ = strconv.Atoi(xvs[1])
		y1, _ = strconv.Atoi(yvs[0])
		y2, _ = strconv.Atoi(yvs[1])
	}

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	maxY := 0
	count := 0
	for xs := -500; xs <= 500; xs++ {
		for ys := -500; ys <= 500; ys++ {
			if v, ok := MaxithinTarget(xs, ys, x1, x2, y1, y2); ok {
				count += 1
				//fmt.Println(xs, ys)
				if v > maxY {
					maxY = v
				}
			}
		}
	}

	fmt.Println("Max Y=", maxY)
	fmt.Println("Count =", count)
}
