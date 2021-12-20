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

var enchancementAlgo = ""

func NewImage() Image {
	return Image{math.MaxInt, math.MaxInt, 0, 0, make(PointMap)}
}

func (img *Image) addPoint(rowNum, colNum int, c rune) {
	if rowNum < img.minRow {
		img.minRow = rowNum
	}
	if rowNum > img.maxRow {
		img.maxRow = rowNum
	}
	if colNum < img.minCol {
		img.minCol = colNum
	}
	if colNum > img.maxCol {
		img.maxCol = colNum
	}
	v := 0
	if c == '#' {
		v = 1
	}
	img.points[Point{rowNum, colNum}] = v
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

func (img Image) Enhance() Image {
	result := NewImage()
	for row := img.minRow - 2; row <= img.maxRow+2; row++ {
		for col := img.minCol - 2; col <= img.maxCol+2; col++ {
			p := img.pixelValue(row, col, 0)
			n := rune(enchancementAlgo[p])
			result.addPoint(row, col, n)
		}
	}
	return result
}

func (img Image) Enhance2(defaultValue int) Image {
	result := NewImage()
	for row := img.minRow - 2; row <= img.maxRow+1; row++ {
		for col := img.minCol - 2; col <= img.maxCol+1; col++ {
			p := img.pixelValue(row, col, defaultValue)
			n := rune(enchancementAlgo[p])
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
		if enchancementAlgo == "" {
			enchancementAlgo = line
		} else {
			inputImage = append(inputImage, line)
		}
	}

	image := NewImage()
	for rowNum, row := range inputImage {
		for colNum, c := range row {
			image.addPoint(rowNum, colNum, c)
		}
	}

	newImage := image.Enhance2(0)
	newImage = newImage.Enhance2(1)
	fmt.Println("Pixels lit =", newImage.CountPixels())

	newImage = image
	for i := 0; i < 50; i++ {
		defaultValue := i % 2
		newImage = newImage.Enhance2(defaultValue)
	}
	fmt.Println("Pixels lit after 50 iterations =", newImage.CountPixels())

	// 29004 is too high
}
