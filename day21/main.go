package main

import (
	"fmt"
	"os"
	"strconv"
)

func Roll(dice int) (int, int) {
	return dice + 1 + dice + 2 + dice + 3, dice + 3
}

func PlayDeterministic(playerPosition1, playerPosition2 int) {
	playerPositions := []int{playerPosition1, playerPosition2}
	playerScores := []int{0, 0}

	currentPlayer := 0
	diceRolls := 0
	dice := 0

	for {
		score, newDice := Roll(dice)
		dice = newDice
		diceRolls += 3

		pos := playerPositions[currentPlayer] - 1
		pos = (pos+score)%10 + 1

		playerPositions[currentPlayer] = pos
		playerScores[currentPlayer] += pos

		if playerScores[currentPlayer] >= 1000 {
			fmt.Println("Player", currentPlayer+1, "wins! Score =", playerScores[currentPlayer])
			break
		}
		currentPlayer = (currentPlayer + 1) % 2
	}

	losingScore := playerScores[(currentPlayer+1)%2]
	completionScore := losingScore * diceRolls
	fmt.Println("Completion Score =", completionScore)
}

type DCacheKey struct {
	player    int
	scores    [2]int // This is a bad cache key
	positions [2]int
}

type DCacheVal [2]uint64

var DCache = make(map[DCacheKey]DCacheVal)

func Dirac(rolls, currentPlayer, dieSum int, scores, positions [2]int, level int) [2]uint64 {
	if rolls < 3 {
		result := [2]uint64{}

		for i := 1; i <= 3; i++ {
			v := Dirac(rolls+1, currentPlayer, dieSum+i, scores, positions, level+1)
			result[0] += v[0]
			result[1] += v[1]
		}

		return result
	} else {
		// One player's turn has ended, update scores and start the next turn
		otherPlayer := (currentPlayer + 1) % 2

		newPos := (positions[currentPlayer]-1+dieSum)%10 + 1
		scores[currentPlayer] = scores[currentPlayer] + newPos
		positions[currentPlayer] = newPos

		key := DCacheKey{currentPlayer, scores, positions}
		if v, ok := DCache[key]; ok {
			//fmt.Println(key, "Using a cached value", v)
			return v
		}

		var v [2]uint64

		if scores[currentPlayer] >= 21 && scores[otherPlayer] < 21 {
			v[currentPlayer] = 1
		} else {
			v = Dirac(0, otherPlayer, 0, scores, positions, level+1)
		}

		//fmt.Println(key, "Storing a cached value")
		DCache[key] = v

		return v
	}
}

func PlayRecursiveDirac(playerPosition1, playerPosition2 int) {
	wins := Dirac(0, 0, 0, [2]int{}, [2]int{playerPosition1, playerPosition2}, 0)
	fmt.Println(wins)
}

func main() {
	playerPosition1 := 4
	playerPosition2 := 8

	if len(os.Args) > 2 {
		playerPosition1, _ = strconv.Atoi(os.Args[1])
		playerPosition2, _ = strconv.Atoi(os.Args[2])
	}

	fmt.Println("Starting positions", playerPosition1, playerPosition2)

	fmt.Println("Part 1")
	PlayDeterministic(playerPosition1, playerPosition2)

	fmt.Println("Part 2")
	PlayRecursiveDirac(playerPosition1, playerPosition2)
}
