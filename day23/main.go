package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
)

const HallwayRow = 1

type Amphipod struct {
	class      rune
	row, col   int
	hasMoved   bool
	energyUsed int
}

func Print(grid [][]rune) {
	for _, row := range grid {
		for _, c := range row {
			fmt.Print(string(c))
		}
		fmt.Println()
	}
}

func (a Amphipod) String() string {
	return fmt.Sprintf("{ %c, row: %d col: %d e: %d}", a.class, a.row, a.col, a.energyUsed)
}

var AmphipodHomes = map[rune]int{'A': 3, 'B': 5, 'C': 7, 'D': 9}
var EnergyCosts = map[rune]int{'A': 1, 'B': 10, 'C': 100, 'D': 1000}

var MinEnergyUsed = math.MaxInt

func addPosition(positions []Amphipod, a Amphipod, col, energySum int, grid [][]rune) []Amphipod {
	energyCost := EnergyCosts[a.class]

	if col == AmphipodHomes[a.class] {
		// We can move into our home room if it is empty or contains only our class
		row := a.row
		for i := HallwayRow + 1; grid[i][col] != '#'; i++ {
			if grid[i][col] != '.' && grid[i][col] != a.class {
				return positions
			}
			if grid[i][col] == '.' {
				row = i
				energySum += energyCost
			}
		}

		positions = append(positions, Amphipod{a.class, row, col, true, a.energyUsed + energySum})
	} else {
		if !a.hasMoved && grid[HallwayRow+1][col] == '#' {
			positions = append(positions, Amphipod{a.class, HallwayRow, col, true, a.energyUsed + energySum + energyCost})
		}
	}

	return positions
}

// If we have a home location in there, then just use that
func trimPositions(positions []Amphipod, a Amphipod) []Amphipod {
	for _, p := range positions {
		if p.col == AmphipodHomes[a.class] {
			return []Amphipod{p}
		}
	}
	return positions
}

func isAtHome(a Amphipod, grid [][]rune) bool {
	if a.col != AmphipodHomes[a.class] {
		return false
	}
	for i := a.row; grid[i][a.col] != '#'; i++ {
		if grid[i][a.col] != a.class {
			return false
		}
	}
	return true
}

func OrganizeAmphipods(amphipods []Amphipod, grid [][]rune) int {
	found := true

	totalEnergyUsed := 0
	for _, a := range amphipods {
		totalEnergyUsed += a.energyUsed
		if a.col != AmphipodHomes[a.class] {
			found = false
		}
	}

	if found {
		if totalEnergyUsed < MinEnergyUsed {
			fmt.Println("Min solution found", totalEnergyUsed)
			MinEnergyUsed = totalEnergyUsed
		}
		return totalEnergyUsed
	}

	if totalEnergyUsed > MinEnergyUsed {
		return totalEnergyUsed
	}

outer:
	for amphipodIndex, a := range amphipods {
		//fmt.Println("Considering amphipod", a)

		if isAtHome(a, grid) {
			continue
		}

		moveToPositions := []Amphipod{}
		energy := 0

		if !a.hasMoved {
			// Still in a room, see if we can step out
			for i := a.row - 1; i >= HallwayRow; i-- {
				if grid[i][a.col] != '.' {
					continue outer
				}
				energy += EnergyCosts[a.class]
			}
		}

		energySum := energy + EnergyCosts[a.class]
		for i := a.col + 1; grid[HallwayRow][i] == '.'; i++ {
			moveToPositions = addPosition(moveToPositions, a, i, energySum, grid)
			energySum += EnergyCosts[a.class]
		}

		energySum = energy
		for i := a.col - 1; grid[HallwayRow][i] == '.'; i-- {
			moveToPositions = addPosition(moveToPositions, a, i, energySum, grid)
			energySum += EnergyCosts[a.class]
		}

		moveToPositions = trimPositions(moveToPositions, a)

		//fmt.Println(moveToPositions)

		for _, newPos := range moveToPositions {
			saveA := a
			a = newPos

			grid[saveA.row][saveA.col] = '.'
			grid[a.row][a.col] = a.class
			amphipods[amphipodIndex] = a
			//Print(grid)

			OrganizeAmphipods(amphipods, grid)

			grid[a.row][a.col] = '.'
			grid[saveA.row][saveA.col] = saveA.class
			amphipods[amphipodIndex] = saveA
			a = saveA
		}
	}

	return totalEnergyUsed
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	grid := make([][]rune, 0)
	amphipods := []Amphipod{}

	rowNum := 0
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		row := make([]rune, len(line))
		for colNum, c := range line {
			if unicode.IsLetter(c) {
				a := Amphipod{c, rowNum, colNum, false, 0}
				amphipods = append(amphipods, a)
			}
			row[colNum] = c
		}
		rowNum += 1
		grid = append(grid, row)
	}

	Print(grid)
	OrganizeAmphipods(amphipods, grid)
	fmt.Println("Energy used =", MinEnergyUsed)
}
