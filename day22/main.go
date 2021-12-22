package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Cuboid struct {
	turnOn                 bool
	x1, x2, y1, y2, z1, z2 int
}

func parseRange(line string) (int, int) {
	line = line[2:]
	ab := strings.Split(line, "..")
	a, _ := strconv.Atoi(ab[0])
	b, _ := strconv.Atoi(ab[1])
	return a, b
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	cuboids := []Cuboid{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		turnOn := true
		if strings.HasPrefix(line, "on") {
			line = strings.TrimPrefix(line, "on ")
		} else {
			line = strings.TrimPrefix(line, "off ")
			turnOn = false
		}

		xyz := strings.Split(line, ",")
		x1, x2 := parseRange(xyz[0])
		y1, y2 := parseRange(xyz[1])
		z1, z2 := parseRange(xyz[2])
		cuboids = append(cuboids, Cuboid{turnOn, x1, x2, y1, y2, z1, z2})
	}

	clamp := func(x int) int {
		if x < -50 {
			return -51
		}
		if x > 50 {
			return 51
		}
		return x
	}

	cubes := make(map[Point]bool)

	for _, c := range cuboids {
		x1 := clamp(c.x1)
		x2 := clamp(c.x2)
		y1 := clamp(c.y1)
		y2 := clamp(c.y2)
		z1 := clamp(c.z1)
		z2 := clamp(c.z2)

		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				for z := z1; z <= z2; z++ {
					if x >= -50 && y >= -50 && z >= -50 && x <= 50 && y <= 50 && z <= 50 {
						cubes[Point{x, y, z}] = c.turnOn
					}
				}
			}
		}

	}

	sum := 0
	for _, v := range cubes {
		if v {
			sum += 1
		}
	}

	fmt.Println(sum)
}
