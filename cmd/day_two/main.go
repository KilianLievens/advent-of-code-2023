package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

func main() {
	maxCubeCounts := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	exampleOneInput := advent.Read("./input/day_two/example_one.txt")
	sum, power := playGame(exampleOneInput, maxCubeCounts)
	fmt.Printf("Example input: sum %d, power %d\n", sum, power) // 8, 2286

	puzzleInput := advent.Read("./input/day_two/puzzle_one.txt")
	sum, power = playGame(puzzleInput, maxCubeCounts)
	fmt.Printf("Puzzle input: sum %d, power %d\n", sum, power) // 2285, 77021
}

func playGame(input []string, maxCounts map[string]int) (int, int) {
	var idSum int
	var powerSum int

	for _, line := range input {
		validGame := true
		observedCounts := map[string]int{}

		gameSegments := strings.Split(line, ": ")
		gameId, _ := strconv.Atoi(strings.ReplaceAll(gameSegments[0], "Game ", ""))
		grabs := strings.Split(gameSegments[1], "; ")

		for _, grab := range grabs {
			colorTexts := strings.Split(grab, ", ")

			for _, colorText := range colorTexts {
				colorSegments := strings.Split(colorText, " ")

				color := colorSegments[1]
				colorCount, _ := strconv.Atoi(colorSegments[0])
				maxColorCount := maxCounts[color]
				oldMax := observedCounts[color]

				if colorCount > maxColorCount {
					validGame = false
				}

				if colorCount > oldMax {
					observedCounts[color] = colorCount
				}
			}
		}

		if validGame {
			idSum += gameId
		}

		gamePower := 1
		for _, count := range observedCounts {
			gamePower *= count
		}

		powerSum += gamePower
	}

	return idSum, powerSum
}
