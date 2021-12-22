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

/*
 * Every throw of the dice generates 3 universes
 * n throws generate thus 3^n universes
 *
 * On the first cast, we get 3 universes
 * Player 1, starting from 4 would get thus 5, 6 and 7 points in the universes 1, 2 and 3 respectivelyj.
 * The second cast happens in universes 1, 2 and 3. There will be 3 outcomes in each of these universes,
 * resulting in 9 universes.
 * Player 2 (starting from 8) would get points 9, 10 and 1 in the 9 universes (9, 10, 1, 9, 10, 1, ...).
 * The second cast happens in 9 universes and results in 27 universes.
 * Player 1 will get points ...
 *
 * Except that each turn the player rolls the dice 3 times and adds up the results
 * So the first turn will create 27 universes that need to be "evaluated" in turn 2
 * It also may be that player 1 wins in some of these 27 universes
 */
func PlayDirac(playerPosition1, playerPosition2 int) {
	wins := [2]int64{}
	playerScores := [2][27]int{}
	playerPositions := [2][27]int{}
	completed := [2]bool{}
	for i := 0; i < 27; i++ {
		playerPositions[0][i] = playerPosition1
		playerPositions[1][i] = playerPosition2
	}

	// We don't actually need all of the 27, there's just 7 unique values here (1+1+1 to 3+3+3)
	// but too lazy
	diceScores := []int{}
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				diceScores = append(diceScores, i+j+k)
			}
		}
	}

	currentPlayer := 0
	universes := int64(1)
	for iteration := 0; !(completed[0] && completed[1]); iteration++ {
		universes *= 27

		// Each universe will split into 27 new universes
		// Track player scores and positions for the 27 possible values

		positions := &playerPositions[currentPlayer]
		scores := &playerScores[currentPlayer]

		fmt.Printf("#%d: uni=%v player%d\n", iteration, universes, currentPlayer+1)
		fmt.Println("scores=", scores)
		fmt.Println("positions=", positions)

		allBigScores := true

		// Run through the results of all the 3 dice throw results in the 27 universes
		for i, diceValue := range diceScores {
			if scores[i] < 21 {
				p := (positions[i]-1+diceValue)%10 + 1
				positions[i] = p
				scores[i] += p

				other := (currentPlayer + 1) % 2
				if scores[i] >= 21 && playerScores[other][i] < 21 {
					wins[currentPlayer] += universes / 27
				} else {
					allBigScores = false
				}
			}
		}

		if allBigScores {
			completed[currentPlayer] = true
			fmt.Printf("Player%d exhausted, %v wins.\n", currentPlayer+1, wins[currentPlayer])
			fmt.Println(scores)
		}

		currentPlayer = (currentPlayer + 1) % 2
	}
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
	PlayDirac(playerPosition1, playerPosition2)
}
