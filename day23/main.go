package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
)

const HallwayRow = 1
const FrontRoomRow = 2
const BackRoomRow = 3

type Amphipod struct {
	class      rune
	row, col   int
	hasMoved   bool
	energyUsed int
}

type Position struct {
	row, col, energy int
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

// Amount of energy required to move to hallway
func RoomToHallwayEnergy(a Amphipod) int {
	return 0
}

var MinEnergyUsed = math.MaxInt

func addPosition(positions []Position, a Amphipod, col, energySum int, grid [][]rune) []Position {
	energyCost := EnergyCosts[a.class]
	homeCol := AmphipodHomes[a.class]

	if col == homeCol {
		if grid[BackRoomRow][col] == '.' {
			positions = append(positions, Position{BackRoomRow, col, energySum + 2*energyCost})
		} else if grid[BackRoomRow][col] == a.class && grid[FrontRoomRow][col] == '.' {
			positions = append(positions, Position{FrontRoomRow, col, energySum + energyCost})
		}
	} else {
		if !a.hasMoved && grid[FrontRoomRow][col] == '#' {
			positions = append(positions, Position{HallwayRow, col, energySum + energyCost})
		}
	}

	return positions
}

// If we have a home location in there, then just use that
func trimPositions(positions []Position, a Amphipod) []Position {
	for _, p := range positions {
		if p.col == AmphipodHomes[a.class] {
			return []Position{p}
		}
	}
	return positions
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

	for amphipodIndex, a := range amphipods {
		//fmt.Println("Considering amphipod", a)
		moveToPositions := []Position{}
		energy := 0

		if !a.hasMoved {
			//fmt.Println("still in a room")
			// Still in a room, see if we can step out
			if a.row == BackRoomRow {
				if grid[FrontRoomRow][a.col] != '.' {
					// Exit blocked
					//fmt.Println("Exit blocked")
					continue
				}
				energy += EnergyCosts[a.class]
			}
			energy += EnergyCosts[a.class]
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

		for _, pos := range moveToPositions {
			saveA := a
			a.hasMoved = true
			a.row = pos.row
			a.col = pos.col
			a.energyUsed = a.energyUsed + pos.energy

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
