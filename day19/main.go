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

type Rotation struct {
	xa, xs, ya, ys, zs int
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

func PointOrientation(p Point, xa, xdir, ya, ydir, zdir int) Point {
	r := Point{}
	r.x = PointVal(p, xa) * xdir
	r.y = PointVal(p, ya) * ydir

	/*
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
	*/

	axis := []int{XAxis, YAxis, ZAxis}
	for _, axis := range axis {
		if xa != axis && ya != axis {
			r.z = PointVal(p, axis) * zdir
		}
	}

	return r
}

func RotateSensors(s Scanner, xa, xd, ya, yd, zs int) Scanner {
	r := NewScanner()
	for _, sv := range s.sensors {
		p := PointOrientation(sv, xa, xd, ya, yd, zs)
		r.addPoints(p.x, p.y, p.z)
	}
	return r
}

func AdjustScanner(s Scanner, correction Point) {
	for i, _ := range s.sensors {
		s.sensors[i].x += correction.x
		s.sensors[i].y += correction.y
		s.sensors[i].z += correction.z
	}
}

func HaveOverlapWith(s1, s2 Scanner, debug bool) (bool, Point) {
	distanceMap := make(map[Point]int)

	for _, v1 := range s1.sensors {
		for _, v2 := range s2.sensors {
			p := Point{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
			if _, ok := distanceMap[p]; ok {
				if debug {
					fmt.Println("Found:", v1, "d=", p)
				}
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

	overlaps := make(map[int]map[int]Point)
	for i := 0; i < len(scanners); i++ {
		overlaps[i] = make(map[int]Point)
	}

	FindOverlappingScanners(0, scanners, make([]bool, len(scanners)), overlaps)
	fmt.Println(overlaps)

	corrections := make([]Point, len(scanners))
	for i := 0; i < len(scanners); i++ {
		_, c := calculateCorrection(i, make(map[int]bool), overlaps)
		corrections[i] = c
		fmt.Println("Correction for", i, "=", c)
	}

	beacons := make(map[Point]int)
	// Scanners are already rotated to align so just sum the beacons up
	for i, s := range scanners {
		AdjustScanner(scanners[i], corrections[i])
		for _, p := range s.sensors {
			beacons[p] += 1
		}
	}

	beaconsFound := len(beacons)

	// 351 is too low
	// 927 is too high
	fmt.Println("Beacons found =", beaconsFound)
}

func FindOverlappingScanners(pos int, scanners []Scanner, found []bool, overlaps map[int]map[int]Point) {
	for i := range scanners {
		if i == 0 || found[i] {
			continue
		}
		if ok, delta, r := FindOverlap(scanners[pos], scanners[i], false); ok {
			fmt.Printf("Overlapping %d and %d: d=%+v r=%+v\n", pos, i, delta, r)
			scanners[i] = RotateSensors(scanners[i], r.xa, r.xs, r.ya, r.ys, r.zs)
			overlaps[pos][i] = delta
			overlaps[i][pos] = delta
			found[i] = true

			FindOverlappingScanners(i, scanners, found, overlaps)
		}
	}
}

func FindOverlap(s1, s2 Scanner, debug bool) (bool, Point, Rotation) {
	axes := []int{XAxis, YAxis, ZAxis}
	signs := []int{1, -1}

	for _, xa := range axes {
		for _, xs := range signs {
			for _, ya := range axes {
				for _, ys := range signs {
					if xa != ya {
						for _, zs := range signs {
							r := RotateSensors(s2, xa, xs, ya, ys, zs)
							if found, delta := HaveOverlapWith(s1, r, false); found {
								if debug {
									fmt.Println("Rotation", xa, xs, ya, ys, "Overlap d=", delta)
								}
								return true, delta, Rotation{xa, xs, ya, ys, zs}
							}
						}
					}
				}
			}
		}
	}
	return false, Point{}, Rotation{}
}

func calculateCorrection(startFrom int, visited map[int]bool, overlaps map[int]map[int]Point) (bool, Point) {
	if startFrom == 0 {
		return true, Point{}
	}
	if haveVisited := visited[startFrom]; haveVisited {
		return false, Point{}
	}

	visited[startFrom] = true
	defer func() {
		visited[startFrom] = false
	}()

	for to, v := range overlaps[startFrom] {
		if found, d := calculateCorrection(to, visited, overlaps); found {
			r := Point{d.x + v.x, d.y + v.y, d.z + v.z}
			return true, r
		}
	}

	return false, Point{}
}
