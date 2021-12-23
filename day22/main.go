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

func (c Cuboid) String() string {
	return fmt.Sprintf("[%v %v %v %v %v %v]", c.x1, c.x2, c.y1, c.y2, c.z1, c.z2)
}

func parseRange(line string) (int, int) {
	line = line[2:]
	ab := strings.Split(line, "..")
	a, _ := strconv.Atoi(ab[0])
	b, _ := strconv.Atoi(ab[1])
	return a, b
}

func clamp(x int) int {
	if x < -50 {
		return -51
	}
	if x > 50 {
		return 51
	}
	return x
}

func CalculatePart1(cuboids []Cuboid) int {
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

	return sum
}

func abs(v int) int {
	if v < 0 {
		return -v
	} else {
		return v
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Volume(c Cuboid) int64 {
	w := abs(c.x2-c.x1) + 1
	h := abs(c.y2-c.y1) + 1
	d := abs(c.z2-c.z1) + 1
	return int64(w * h * d)
}

func Contains(c Cuboid, x, y, z int) bool {
	return x >= c.x1 && x <= c.x2 && y >= c.y1 && y <= c.y2 && z >= c.z1 && z <= c.z2
}

// Returns the cuboid containing the overlap of c1 and c2
func Overlap(c1, c2 Cuboid, sign bool) (Cuboid, bool) {
	// c1 fully contains c2
	if c1.x1 <= c2.x1 && c1.x2 >= c2.x2 && c1.y1 <= c2.y1 && c1.y2 >= c2.y2 && c1.z1 <= c2.z1 && c1.z2 >= c2.z2 {
		c2.turnOn = sign
		return c2, true
	}
	// c2 fully contains c1
	if c2.x1 <= c1.x1 && c2.x2 >= c1.x2 && c2.y1 <= c1.y1 && c2.y2 >= c1.y2 && c2.z1 <= c1.z1 && c2.z2 >= c1.z2 {
		c1.turnOn = sign
		return c1, true
	}

	if c2.x2 >= c1.x1 && c2.x1 <= c1.x2 && c2.y2 >= c1.y1 && c2.y1 <= c1.y2 && c2.z2 >= c1.z1 && c2.z1 <= c1.z2 {
		cx1 := max(c2.x1, c1.x1)
		cx2 := min(c2.x2, c1.x2)

		cy1 := max(c2.y1, c1.y1)
		cy2 := min(c2.y2, c1.y2)

		cz1 := max(c2.z1, c1.z1)
		cz2 := min(c2.z2, c1.z2)

		return Cuboid{sign, cx1, cx2, cy1, cy2, cz1, cz2}, true
	}
	return Cuboid{}, false
}

func CountCubesLit(cuboids []Cuboid, level int) int64 {
	left := []Cuboid{}
	right := cuboids

	cubesLit := int64(0)

	for len(right) > 0 {
		rightCuboid := right[0]
		right = right[1:]

		overlaps := []Cuboid{}
		for _, leftCuboid := range left {
			if overlap, found := Overlap(leftCuboid, rightCuboid, leftCuboid.turnOn); found {
				overlaps = append(overlaps, overlap)
			}
		}

		left = append(left, rightCuboid)
		cubesFound := Volume(rightCuboid)

		if rightCuboid.turnOn {
			litOverlaps := CountCubesLit(overlaps, level+1)
			cubesLit += cubesFound - litOverlaps
		} else {
			turnOffOverlaps := CountCubesLit(overlaps, level+1)
			cubesLit -= turnOffOverlaps
		}
	}

	return cubesLit
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

	sum := CalculatePart1(cuboids)
	fmt.Println("Part 1 =", sum)

	cubesLit := CountCubesLit(cuboids, 0)
	fmt.Println("Part 2 =", cubesLit)
}
