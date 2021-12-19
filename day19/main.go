package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const XAxis = 0
const YAxis = 1
const ZAxis = 2

type Point struct {
	x, y, z int
}

type SensorMap map[Point]bool

type Scanner struct {
	sensors []Point
}

func NewScanner() Scanner {
	return Scanner{
		sensors: make([]Point, 0),
	}
}

func (s *Scanner) Print() {
	fmt.Println("--- scanner ---")
	for _, s := range s.sensors {
		fmt.Printf("%d,%d,%d\n", s.x, s.y, s.z)
	}
	fmt.Println()
}

func (s *Scanner) addPoints(x, y, z int) {
	s.sensors = append(s.sensors, Point{x, y, z})
}

func PointVal(p Point, axis int) int {
	switch axis {
	case XAxis:
		return p.x
	case YAxis:
		return p.y
	default:
		return p.z
	}
}

func PointOrientation(p Point, xa, xdir, ya, ydir int) Point {
	r := Point{}
	r.x = PointVal(p, xa) * xdir
	r.y = PointVal(p, ya) * ydir

	zdir := 1
	if xa == XAxis && ya == ZAxis {
		zdir = -1
	}
	if xa == YAxis && ya == XAxis {
		zdir = -1
	}
	if xa == ZAxis && ya == YAxis {
		zdir = -1
	}

	zdir *= xdir * ydir

	axis := []int{XAxis, YAxis, ZAxis}
	for _, axis := range axis {
		if xa != axis && ya != axis {
			r.z = PointVal(p, axis) * zdir
		}
	}

	return r
}

func RotateSensors(s Scanner, xa, xd, ya, yd int) Scanner {
	r := NewScanner()
	for _, sv := range s.sensors {
		p := PointOrientation(sv, xa, xd, ya, yd)
		r.addPoints(p.x, p.y, p.z)
	}
	return r
}

func HaveOverlapWith(s1, s2 Scanner) (bool, Point) {
	distanceMap := make(map[Point]int)

	for _, v1 := range s1.sensors {
		for _, v2 := range s2.sensors {
			p := Point{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
			if _, ok := distanceMap[p]; ok {
				//fmt.Println("Found:", v1, "d=", p)
			}
			distanceMap[p] += 1
			if distanceMap[p] >= 12 {
				return true, p
			}
		}
	}

	return false, Point{}
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	currentScanner := -1
	scanners := []Scanner{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "---") {
			currentScanner += 1
			scanners = append(scanners, NewScanner())
		} else {
			coords := strings.Split(line, ",")
			x, _ := strconv.Atoi(coords[0])
			y, _ := strconv.Atoi(coords[1])
			z, _ := strconv.Atoi(coords[2])
			scanners[currentScanner].addPoints(x, y, z)
		}
	}

	fmt.Println(len(scanners), "scanners loaded")

	for i := 0; i < len(scanners)-1; i++ {
		for j := i + 1; j < len(scanners); j++ {
			fmt.Printf("Findinig overlaps for scanners %d and %d\n", i, j)
			FindOverlap(scanners[i], scanners[j])
		}
	}

}

func FindOverlap(s1, s2 Scanner) bool {
	axes := []int{XAxis, YAxis, ZAxis}
	signs := []int{1, -1}

	for _, xa := range axes {
		for _, xs := range signs {
			for _, ya := range axes {
				for _, ys := range signs {
					if xa != ya {
						//fmt.Printf("Rotation xa=%d xs=%d ya=%d ys=%d\n", xa, xs, ya, ys)
						r := RotateSensors(s2, xa, xs, ya, ys)
						if found, delta := HaveOverlapWith(s1, r); found {
							fmt.Println("Overlap found. d=", delta)
							return true
						}
					}
				}
			}
		}
	}
	return false
}
