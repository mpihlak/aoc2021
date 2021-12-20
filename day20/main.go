package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	row, col int
}

type PointMap map[Point]int

type Image struct {
	minRow, minCol, maxRow, maxCol int
	points                         PointMap
}

var EnchancementAlgo = make([]int, 0)

func NewImage() Image {
	return Image{math.MaxInt, math.MaxInt, 0, 0, make(PointMap)}
}

func (img *Image) addPoint(row, col int, v int) {
	if row < img.minRow {
		img.minRow = row
	}
	if row > img.maxRow {
		img.maxRow = row
	}
	if col < img.minCol {
		img.minCol = col
	}
	if col > img.maxCol {
		img.maxCol = col
	}
	img.points[Point{row, col}] = v
}

func (img Image) pixelValue(row, col, defaultValue int) int {
	v := 0
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			v = v << 1
			if pixel, ok := img.points[Point{i, j}]; ok {
				v = v | pixel
			} else {
				v = v | defaultValue
			}
		}
	}
	return v
}

func (img Image) Enhance(defaultValue int) Image {
	result := NewImage()
	for row := img.minRow - 2; row <= img.maxRow+1; row++ {
		for col := img.minCol - 2; col <= img.maxCol+1; col++ {
			p := img.pixelValue(row, col, defaultValue)
			n := EnchancementAlgo[p]
			result.addPoint(row, col, n)
		}
	}
	return result
}

func (img Image) Print() {
	fmt.Printf("Dimensions minRow: %d maxRow: %d, minCol: %d, maxCol: %d\n", img.minRow, img.maxRow, img.minCol, img.maxCol)
	for row := img.minRow; row <= img.maxRow; row++ {
		for col := img.minCol; col <= img.maxCol; col++ {
			p := img.points[Point{row, col}]
			if p > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (img Image) CountPixels() int {
	sum := 0
	for _, v := range img.points {
		sum += v
	}
	return sum
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	inputImage := []string{}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		if len(EnchancementAlgo) == 0 {
			for _, c := range line {
				v := 0
				if c == '#' {
					v = 1
				}
				EnchancementAlgo = append(EnchancementAlgo, v)
			}
		} else {
			inputImage = append(inputImage, line)
		}
	}

	image := NewImage()
	for rowNum, row := range inputImage {
		for colNum, c := range row {
			v := 0
			if c == '#' {
				v = 1
			}
			image.addPoint(rowNum, colNum, v)
		}
	}

	newImage := image.Enhance(0)
	newImage = newImage.Enhance(1)
	fmt.Println("Pixels lit =", newImage.CountPixels())

	newImage = image
	for i := 0; i < 50; i++ {
		defaultValue := i % 2
		newImage = newImage.Enhance(defaultValue)
	}
	fmt.Println("Pixels lit after 50 iterations =", newImage.CountPixels())
}
