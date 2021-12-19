package main

import (
	"fmt"
	"math"
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
	sensors          []Point
	minX, minY, minZ int
}

func NewScanner() Scanner {
	return Scanner{
		sensors: make([]Point, 0),
		minX:    math.MaxInt,
		minY:    math.MaxInt,
		minZ:    math.MaxInt,
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
	if x < s.minX {
		s.minX = x
	}
	if y < s.minY {
		s.minY = y
	}
	if z < s.minZ {
		s.minZ = z
	}
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

	//fmt.Printf("xa=%d xd=%d ya=%d yd=%d\t %v,%v,%v\n", xa, xdir, ya, ydir, r.x, r.y, r.z)

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

func (s *Scanner) Normalize() {
	for k := range s.sensors {
		p := &s.sensors[k]
		p.x -= s.minX
		p.y -= s.minY
		p.z -= s.minZ
	}
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

	//p := Point{-2, -3, 1}

	axes := []int{XAxis, YAxis, ZAxis}
	signs := []int{1, -1}
	for _, xa := range axes {
		for _, xs := range signs {
			for _, ya := range axes {
				for _, ys := range signs {
					if xa != ya {
						s1 := scanners[0]
						s2 := RotateSensors(s1, xa, xs, ya, ys)
						s2.Print()
					}
				}
			}
		}
	}
}
